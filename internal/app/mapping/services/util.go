package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/uibricks/studio-engine/internal/app/mapping/constants"
	"github.com/uibricks/studio-engine/internal/app/mapping/utils"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	expressionpb "github.com/uibricks/studio-engine/internal/pkg/proto/expression"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"github.com/uibricks/studio-engine/internal/pkg/redis"
	globalutils "github.com/uibricks/studio-engine/internal/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (m *MappingServer) GetMappingFromCache(ctx context.Context, projectId string) (*mappingpb.Repositories, error) {
	val, err := m.Redis.HGet(ctx, "test", redis.Misc, "admin", projectId)
	if !redis.IsNil(err) {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch mapping from cache for project id(%s), %v", projectId, err))
	}
	return utils.CastStringToMapping(val)
}

func (m *MappingServer) SetMappingInCache(ctx context.Context, mapping *mappingpb.Repositories, projectId string) error {
	// if this is new entry update the project id and also update created time stamp
	if len(mapping.ProjectId) == 0 {
		mapping.ProjectId = projectId
		mapping.CreatedAt = globalutils.DbFormatDate(time.Now())
	}

	res, err := utils.MarshalToString(mapping)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to marshal mapping to string for the project %s, %v", projectId, err))
	}

	err = m.Redis.HSet(ctx, "test", redis.Misc, "admin", mapping.GetProjectId(), res)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to update the mapping in cache with project %s, %v", projectId, err))
	}

	return nil
}

// LoadMappingWithVersion will load the mapping from db if the project and project version passed.
func (m *MappingServer) LoadMappingWithVersion(mapping *mappingpb.Repositories, projectId string, projectVersion int32) error {
	if projectVersion > 0 {
		err := m.DB.Model(mapping).Where("project_id = ? and project_version = ? and deleted_at is null", projectId, projectVersion).First()
		if err == pg.ErrNoRows {
			return status.Error(codes.NotFound, fmt.Sprintf("Mapping not found with the project(%s) and version(%d)", projectId, projectVersion))
		}
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("Not able to read the config for the project %s and version %d - %v", projectId, projectVersion, err))
		}
	}
	return nil
}

func (m *MappingServer) MoveReposToHistory(ctx context.Context, projectId string, projectVersion int32) {

	tx, err := m.DB.Begin()
	defer tx.Close()

	projectVersionClause := fmt.Sprintf(" and %s < %d", constants.MappingTableProjectColumnName, projectVersion)
	_, err = m.DB.Query(nil, fmt.Sprintf(constants.QueryCopyRepositoriesToHistory, projectId, projectVersionClause))

	if err != nil {
		logger.WithContext(ctx).Errorf("failed to insert data into repositories history table, internal error, %v", err)
		return
	}

	_, err = m.DB.Query(nil, fmt.Sprintf(constants.QueryDeleteObsoleteRepositories, projectId, projectVersionClause))

	if err != nil {
		_ = tx.Rollback()
		logger.WithContext(ctx).Errorf("failed to delete obsolete data from repositories table, internal error, %v", err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithContext(ctx).Errorf("failed to commit the moving repos to history transaction, %v", err)
	}

}

func (m *MappingServer) SoftDeleteMapping(ctx context.Context, req *mappingpb.DeleteMappingRequest) (bool, error) {
	// check in cache
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())
	if err != nil {
		return false, err
	}

	tx, err := m.DB.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Close()

	mapping.ProjectId = req.GetProjectId()
	mapping.Id = 0
	mapping.ProjectVersion = req.GetCachedProjectVersion()
	_, err = m.DB.Model(mapping).Insert()
	if err != nil {
		return false, status.Error(codes.Internal, fmt.Sprintf("Failed to save cached mapping as deleted with project id(%s), %v", req.GetProjectId(), err))
	}

	rows, err := m.DB.Model(&mappingpb.Repositories{}).Set("deleted_at = current_timestamp").Where("project_id=? and deleted_at is null", req.GetProjectId()).Update()
	if err != nil {
		tx.Rollback()
		return false, status.Error(codes.Internal, fmt.Sprintf("failed to delete mapping with project id(%s), %v", req.GetProjectId(), err))
	}

	// also remove from cache
	if rows.RowsAffected() > 0 {
		if len(mapping.GetProjectId()) > 0 {
			err = m.Redis.HDel(ctx, "test", redis.Misc, "admin", req.GetProjectId())
			if err != nil {
				return false, status.Error(codes.Internal, fmt.Sprintf("Failed to delete-mapping from cache for project-id : %s, %v", req.GetProjectId(), err))
			}
		}
		return true, nil
	}
	return false, nil
}

func (m *MappingServer) PermanentlyDeleteMapping(ctx context.Context, projectId string) (bool, error) {

	// delete mapping permanently
	rows, err := m.DB.Model(&mappingpb.Repositories{}).Where("project_id=?", projectId).Delete()
	if err != nil {
		errMsg := fmt.Sprintf("failed to deleted mapping permanently, internal error, %v", err)
		logger.WithContext(ctx).Error(errMsg)
		return false, fmt.Errorf(errMsg)
	}

	// if mapping is deleted then
	if rows.RowsAffected() > 0 {
		// delete mapping from history
		_, err = m.DB.Model(&mappingpb.Repositories{}).Exec(constants.QueryDeleteFromRepositoriesHistory, projectId)
		if err != nil {
			errMsg := fmt.Sprintf("failed to deleted mapping permanently from history table for project-id(%s), internal error, %v", projectId, err)
			logger.WithContext(ctx).Error(errMsg)
		}
		return true, nil
	}
	return false, nil
}

func (m *MappingServer) ExecuteAction(ctx context.Context, rmqMsg *rabbitmq.RMQMessage) error {

	switch rmqMsg.Action {
	case rabbitmq.Action_Save_Mapping:

		req := &mappingpb.SaveMappingRequest{}
		b, err := json.Marshal(rmqMsg.Payload)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to marshal payload to string payload(%v),%v", rmqMsg.Payload, err))
		}

		json.Unmarshal(b, req)

		_, err = m.SaveMapping(ctx, req)
		return err

	case rabbitmq.Action_Delete_Mapping:

		req := &mappingpb.DeleteMappingRequest{}
		b, err := json.Marshal(rmqMsg.Payload)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to marshal payload to string payload(%v),%v", rmqMsg.Payload, err))
		}

		json.Unmarshal(b, req)

		_, err = m.DeleteMapping(ctx, req)

		return err
	case rabbitmq.Action_Restore_Mapping:
		req := &mappingpb.RestoreMappingRequest{}
		b, err := json.Marshal(rmqMsg.Payload)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to marshal payload to string payload(%v),%v", rmqMsg.Payload, err))
		}

		json.Unmarshal(b, req)

		_, err = m.RestoreMapping(ctx, req)

		return err
	default:
		return fmt.Errorf("Invalid action sent in mapping queue : %s", rmqMsg.Action)
	}

}

func PopulateGroupChildren(expressions *[]*expressionpb.Menu, expMenu []*mappingpb.Menu) {

	for _,e := range *expressions {
		if e.GetType() == constants.GroupExpressionType {
			children := GetExpressionChildren(expMenu, e.GetId())
			if children != nil {
				e.Children = CastMappingMenuToExpressionMenu(children)
			}
		}
	}
}

func GetExpressionChildren(expMenu []*mappingpb.Menu, expId string) []*mappingpb.Menu {
	for _,e := range expMenu {
		if e.Id == expId {
			return e.Children
		}
	}
	return nil
}

func CastMappingMenuToExpressionMenu(menu []*mappingpb.Menu) []*expressionpb.Menu {
	expMenu := make([]*expressionpb.Menu, 0)
	b, _ := json.Marshal(menu)
	_ = json.Unmarshal(b, &expMenu)
	return expMenu
}