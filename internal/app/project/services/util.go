package project

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/uibricks/studio-engine/internal/app/project/constants"
	"github.com/uibricks/studio-engine/internal/app/project/utils"
	_ "github.com/uibricks/studio-engine/internal/pkg/constants"
	pkgConstants "github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	projectpb "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	"github.com/uibricks/studio-engine/internal/pkg/redis"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	pkgUtils "github.com/uibricks/studio-engine/internal/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func (p *ProjectServer) GetProjectFromCache(ctx context.Context, projectId string) (*projectpb.Object, error) {
	val, err := p.Redis.HGet(ctx, "test", redis.Misc, "admin", projectId)
	if !redis.IsNil(err) {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to fetch project with id : %s from cache, %v", projectId, err))
	}
	return utils.CastStringToProject(val)
}

func (p *ProjectServer) SetProjectInCache(ctx context.Context, project *projectpb.Object) error {
	res, err := utils.MarshalToString(project)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to marshal project to string with id : %s, %v", project.GetLuid(), err))
	}

	err = p.Redis.HSet(ctx, "test", redis.Misc, "admin", project.GetLuid(), res)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to set the project in cache, %v", err))
	}

	return nil
}

func (p *ProjectServer) GetProjectsFromCache(ctx context.Context) ([]*projectpb.Object, error) {
	val, err := p.Redis.HGetAll(ctx, "test", redis.Misc, "admin")
	if !redis.IsNil(err) {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch all projects from cache, %v", err))
	}
	return utils.CastStringToProjects(val)
}

func (p *ProjectServer) UpdateProjectStatus(ctx context.Context, projectStatus string, project *projectpb.Object) error {
	if projectStatus != pkgConstants.STATE_ACTIVE {
		projectStatus = projectStatus + "-" + request.GetContextRequestID(ctx)
	}

	// to send updated project state back to gateway [instead of initially set 'pending' state]
	project.State = projectStatus

	_, err := p.DB.Model(project).Set("state = '"+projectStatus+"', updated_at = ?", pkgUtils.DbFormatDate(time.Now())).Where("id = ?", project.GetId()).Update()
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to update the project status to %s, internal error, %v", projectStatus, err))
	}
	return nil
}

func (p *ProjectServer) RecoverSaveProject(ctx context.Context, project *projectpb.Object, panicErr interface{}) (*projectpb.SaveProjectResponse, error) {
	logger.WithContext(ctx).Errorf("Panic in saving project saga - %v.", panicErr)
	savedProject := &projectpb.Object{}
	err := p.DB.Model(savedProject).Where("luid::text=?0 and user_version=?1", project.GetLuid(), project.UserVersion).First()
	if err == pg.ErrNoRows {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to save the project please contact administrator project -(%s)", project.GetLuid()))
	} else {
		if savedProject.GetState() == pkgConstants.STATE_ACTIVE {
			return &projectpb.SaveProjectResponse{UpdatedAt: savedProject.GetUpdatedAt(), Luid: savedProject.Luid, Status: "Project saved successfully"}, nil
		}
		if err := p.UpdateProjectStatus(ctx, pkgConstants.STATE_PANIC, project); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to save the project -(%s)", project.GetLuid()))
		}
	}
	return nil, status.Error(codes.Internal, fmt.Sprintf("failed to save the project -(%s)", project.GetLuid()))
}

func (p *ProjectServer) GetCurrentActiveProject(projectId string) (*projectpb.Object, error) {
	project := &projectpb.Object{}
	err := p.DB.Model(project).Where("luid::text=?0 and deleted_at is null and state='"+pkgConstants.STATE_ACTIVE+"'", projectId).Order("updated_at desc").First()
	return project, err
}

func (p *ProjectServer) MoveProjectsToHistory(ctx context.Context, luid string) {

	tx, err := p.DB.Begin()
	defer tx.Close()

	_, err = p.DB.Query(nil, fmt.Sprintf(constants.QueryCopyProjectsToHistory, luid))

	if err != nil {
		logger.WithContext(ctx).Errorf("failed to insert data into project history table, internal error, %v", err)
		return
	}

	_, err = p.DB.Query(nil, fmt.Sprintf(constants.QueryDeleteObsoleteProjects, luid))

	if err != nil {
		_ = tx.Rollback()
		logger.WithContext(ctx).Errorf("failed to delete obselete data from objects table, %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		logger.WithContext(ctx).Errorf("failed to commit the moving projects to history transaction, %v", err)
	}
}

func (p *ProjectServer) SoftDeleteProject(ctx context.Context, luid string) (bool, error) {
	// get the project from cache
	project, err := p.GetProjectFromCache(ctx, luid)
	if err != nil {
		return false, err
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Close()

	// TODO - We need to cache the projects for all the users when we integrate authentication
	// if project exists in cache then save the cached data into db
	if len(project.GetLuid()) > 0 {

		project.Id = 0
		project.Version = 0
		project.State = pkgConstants.STATE_CACHE

		// save cached project to db with status as 'cached'
		// to indicate the project was never saved but got saved to db from cache during delete
		_, err = p.DB.Model(project).Insert()
		if err != nil {
			return false, status.Error(codes.Internal, fmt.Sprintf("Failed to insert project to db as deleted, %v", err))
		}
	}

	res, err := p.DB.Model(&projectpb.Object{}).Set("deleted_at=current_timestamp").Where("deleted_at is null and luid::text=?", luid).Update()

	if err != nil {
		tx.Rollback()
		return false, status.Error(codes.Internal, fmt.Sprintf("Failed to delete the project in db, %v", err))
	}

	if res.RowsAffected() > 0 {
		// removing from cache as well
		if len(project.GetLuid()) > 0 {
			err = p.Redis.HDel(ctx, "test", redis.Misc, "admin", luid)
			if err != nil {
				logger.WithContext(ctx).Errorf("Failed to delete-project from cache for project-id : %s, %v", luid, err)
			}
		}
		return true, nil
	}

	return false, nil
}

func (p *ProjectServer) PermanentlyDeleteProject(ctx context.Context, luid string) (bool, error) {

	// delete the project permanently
	res, err := p.DB.Model(&projectpb.Object{}).Where("luid::text=?", luid).Delete()

	if err != nil {
		return false, status.Error(codes.Internal, fmt.Sprintf("failed to permanently delete the project with id (%s), err - %v", luid, err))
	}

	// if the project is deleted then
	if res.RowsAffected() > 0 {
		// delete the project from history
		_, err = p.DB.Model(&projectpb.Object{}).Exec(constants.QueryDeleteFromProjectHistory, luid)

		// if it fails to delete from history and log and dont return error because
		// it has been deleted from main table. User can never use this project.
		if err != nil {
			logger.WithContext(ctx).Errorf("failed to delete the project (%s) history permanently, err - %v", luid, err)
		}

		return true, nil
	}

	return false, nil
}

func (p *ProjectServer) GetProjectsCount(luid string) (int32, error) {

	count, err := p.DB.Model(&projectpb.Object{}).Where(fmt.Sprintf("luid::text='%s' and deleted_at is null and state='%s'", luid, pkgConstants.STATE_ACTIVE)).Count()
	if err != nil {
		return -1, status.Error(codes.Internal, fmt.Sprintf("Failed to fetch data from db for getting projects count, internal error, %v", err))
	}

	return int32(count), nil
}

func (p *ProjectServer) GetAllProjects(ctx context.Context, parentId int32) ([]*projectpb.Object, error) {

	var rows []*projectpb.Object

	// get all projects from cache
	projects, err := p.GetProjectsFromCache(ctx)

	if err != nil {
		return nil, err
	}

	// remove cached-projects with unmatch parentId
	for i := len(projects) - 1; i >= 0; i-- {
		cacheParentId, _ := strconv.Atoi(projects[i].GetParentId())
		projects[i].State = pkgConstants.STATE_CACHE
		if int32(cacheParentId) != parentId {
			projects = append(projects[:i], projects[i+1:]...)
		}
	}

	// Get projects from DB as well
	_, err = p.DB.Query(&rows, fmt.Sprintf(constants.QueryFetchActiveProjects, utils.GetParentIDFilter(parentId)))

	if err != nil {
		return nil, err
	}

	rows = append(rows, projects...)
	rows = utils.RetainLatestProjects(rows)

	// return all projects in cache and db
	return rows, nil
}

func (p *ProjectServer) GetTrashedProjects(ctx context.Context) ([]*projectpb.Object, error) {

	var rows []*projectpb.Object

	//err := p.DB.Model(&rows).Where(constants.QueryFetchTrashProjects).Select()
	_, err := p.DB.Query(&rows, constants.QueryFetchTrashProjects)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to fetch trashed projects from db, %v", err))
	}

	return rows, nil

}

func (p *ProjectServer) GetRecentProjects(ctx context.Context) ([]*projectpb.Object, error) {

	// this can be changed
	recentDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	rows := make([]*projectpb.Object, 0)

	// get all projects from cache
	projects, err := p.GetProjectsFromCache(ctx)

	if err != nil {
		return nil, err
	}

	// get projects from db greater than given date
	_, err = p.DB.Query(&rows, fmt.Sprintf(constants.QueryFetchRecentProjects, recentDate))

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get the recent projects from db, %v", err))
	}

	// filter cached projects by greater than given date
	for i, project := range projects {
		project.State = pkgConstants.STATE_CACHE
		if project.GetUpdatedAt() > recentDate {
			rows = append(rows, projects[i])
		}
	}

	rows = utils.RetainLatestProjects(rows)

	return rows, nil

}
