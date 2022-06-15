package project

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/streadway/amqp"
	"github.com/uibricks/studio-engine/internal/app/project/constants"
	"github.com/uibricks/studio-engine/internal/app/project/utils"
	pkgConstants "github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/db"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	projectpb "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"github.com/uibricks/studio-engine/internal/pkg/redis"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	pkgUtils "github.com/uibricks/studio-engine/internal/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type ProjectServer struct {
	DB         *pg.DB
	Redis      *redis.Redis
	Rabbit     *rabbitmq.Rabbit
	ReplyQueue amqp.Queue
	Channel    *rabbitmq.Channel
}

func ProvideProjectServer(dbClient db.DbClient, redis *redis.Redis, rabbit *rabbitmq.Rabbit, replyQ amqp.Queue, ch *rabbitmq.Channel) *ProjectServer {
	return &ProjectServer{
		DB:         dbClient.Connection,
		Redis:      redis,
		Rabbit:     rabbit,
		ReplyQueue: replyQ,
		Channel:    ch,
	}
}

func (p *ProjectServer) UpdateProject(ctx context.Context, req *projectpb.ProjectRequest) (*projectpb.ProjectResponse, error) {

	// get the project from cache
	c, err := p.GetProjectFromCache(ctx, req.GetProject().GetLuid())
	if err != nil {
		return nil, err
	}

	// if the cache doesn't exist and this is new project request then make it as new insert
	if c.GetLuid() == "" && req.GetNewProject() {
		c.Luid = req.GetProject().Luid
		c.CreatedAt = pkgUtils.DbFormatDate(time.Now())
	}

	// If not in cache, first fetch latest active from DB
	// If no active projects available, fetch latest one
	if len(c.GetLuid()) == 0 {
		c, err = p.GetCurrentActiveProject(req.GetProject().Luid)

		if err == pg.ErrNoRows {
			return nil, utils.ProjectNotFoundError(req.GetProject().GetLuid())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch project with id : %s from db, %v", req.GetProject().GetLuid(), err))
		}
	}

	// update project object
	if len(req.GetProject().GetName()) > 0 {
		c.Name = req.GetProject().Name
	}

	if req.GetProject().GetConfig() != nil {
		if req.GetProject().GetConfig().GetComponents() != nil {
			if c.GetConfig() == nil {
				c.Config = &projectpb.Config{}
			}
			if c.GetConfig().GetComponents() == nil {
				c.GetConfig().Components = make(map[string]*projectpb.Component)
			}
			for key, comp := range req.GetProject().GetConfig().GetComponents() {
				if comp.GetCreatedAt() == "" {
					if c.GetConfig().Components[key] != nil || c.GetConfig().Components[key].GetCreatedAt() != "" {
						comp.CreatedAt = c.GetConfig().Components[key].GetCreatedAt()
					} else {
						comp.CreatedAt = pkgUtils.DbFormatDate(time.Now())
					}
				}
				comp.UpdatedAt = pkgUtils.DbFormatDate(time.Now())
				c.GetConfig().Components[key] = comp
			}
		}
	}

	c.UpdatedAt = pkgUtils.DbFormatDate(time.Now())

	// set updated project back to cache
	err = p.SetProjectInCache(ctx, c)

	if err != nil {
		return nil, err
	}

	return &projectpb.ProjectResponse{Project: c}, nil
}

func (p *ProjectServer) SaveProject(ctx context.Context, req *projectpb.SaveProjectRequest) (resp *projectpb.SaveProjectResponse, respErr error) {

	// get the project from cache
	project, err := p.GetProjectFromCache(ctx, req.GetLuid())
	if err != nil {
		return nil, err
	}

	if len(project.GetLuid()) == 0 {
		return nil, utils.ProjectNotFoundError(req.GetLuid())
	}

	// to revoke changes in DB in case of panic
	defer func() {
		if panicErr := recover(); panicErr != nil {
			resp, respErr = p.RecoverSaveProject(ctx, project, panicErr)
		}
	}()

	// This lets DB to self increment ProjectID and ProjectVersion in case they are set in cache
	// Prevents Unique/Primary constraint violation
	project.Id = 0
	project.Version = 0

	err = p.SaveProjectSaga(ctx, req, project)

	if err != nil {
		logger.WithContext(ctx).Errorf(fmt.Sprintf("Failed to save the project with id - %s, error - %v", req.GetLuid(), err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to save the project with id - %s, error - %v", req.GetLuid(), err))
	} else {
		p.MoveProjectsToHistory(ctx, req.GetLuid())
	}

	// remove from cache, so that saved project from db is set to cache when update or get project is invoked
	err = p.Redis.HDel(ctx, "test", redis.Misc, "admin", req.GetLuid())
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to delete the cache for the project (%s)", req.GetLuid())
	}

	return &projectpb.SaveProjectResponse{Luid: req.GetLuid(), Status: "Project Saved with version : " + project.UserVersion + " as '" + project.GetState() + "'", UpdatedAt: project.GetUpdatedAt()}, nil
}

func (p *ProjectServer) GetProject(ctx context.Context, req *projectpb.ProjectLuidRequest) (*projectpb.ProjectResponse, error) {
	// get the project from cache
	project, err := p.GetProjectFromCache(ctx, req.GetLuid())

	if err != nil {
		return nil, err
	}

	// If not in cache, first fetch latest active from DB
	if len(project.GetLuid()) == 0 {

		err = p.DB.Model(project).Where("luid::text=?0 and deleted_at is null and state='active'",
			req.GetLuid()).Order("updated_at desc").First()

		if err == pg.ErrNoRows {
			return nil, utils.ProjectNotFoundError(req.GetLuid())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch project with id : %s from db while getting project, internal error, %v", req.GetLuid(), err))
		}

		// set project to cache
		err = p.SetProjectInCache(ctx, project)

		if err != nil {
			return nil, err
		}
	}

	return &projectpb.ProjectResponse{Project: project}, nil
}

func (p *ProjectServer) RestoreProject(ctx context.Context, req *projectpb.ProjectLuidRequest) (resp *projectpb.RestoreProjectResponse, respErr error) {

	err := p.DB.Model(&projectpb.Object{}).Where("luid::text=?0 and deleted_at is null and (state=?1 or state=?2)",
		req.GetLuid(), pkgConstants.STATE_ACTIVE, pkgConstants.STATE_CACHE).Select()

	if err == nil || err != pg.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("cannot restore. non-deleted project(s) found with id : %s", req.GetLuid()))
	}

	err = p.RestoreProjectSaga(ctx, req)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to restore the project with id - %s, error - %v", req.GetLuid(), err))
	}
	return &projectpb.RestoreProjectResponse{UpdatedAt: pkgUtils.DbFormatDate(time.Now()), Status: "Project restored successfully."}, nil
}

func (p *ProjectServer) GetProjects(ctx context.Context, req *projectpb.ProjectsRequest) (*projectpb.ProjectsResponse, error) {

	var projects []*projectpb.Object
	var err error

	switch strings.ToLower(req.GetType()) {

	case pkgConstants.ProjectTypeTrash:
		projects, err = p.GetTrashedProjects(ctx)
		if err != nil {
			return nil, err
		}
		break
	case pkgConstants.ProjectTypeRecent:
		projects, err = p.GetRecentProjects(ctx)
		if err != nil {
			return nil, err
		}
		break
	default:
		projects, err = p.GetAllProjects(ctx, req.GetParentId())
		if err != nil {
			return nil, err
		}
	}

	return &projectpb.ProjectsResponse{Projects: projects}, nil
}

func (p *ProjectServer) DeleteProject(ctx context.Context, req *projectpb.ProjectLuidRequest) (*projectpb.DeleteProjectResponse, error) {

	corrId := request.GetContextRequestID(ctx)

	softDeleted, err := p.SoftDeleteProject(ctx, req.GetLuid())

	if err != nil {
		return nil, err
	}

	permDeleted := false

	if !softDeleted {
		permDeleted, err = p.PermanentlyDeleteProject(ctx, req.GetLuid())

		if err != nil {
			return nil, err
		}
	}

	if !softDeleted && !permDeleted {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("project not found with id %s", req.GetLuid()))
	}

	// prepare payload to delete repos associated with project, permanently
	b, _ := p.Rabbit.PrepareRMQMessage(ctx,
		rabbitmq.Action_Delete_Mapping, &mappingpb.DeleteMappingRequest{
			ProjectId: req.GetLuid(),
		})

	// publish to mapping queue
	err = p.Rabbit.Publish(rabbitmq.Mapping_Queue_Name, corrId, b, p.Channel)

	if err != nil {
		logger.WithContext(ctx).Errorf("failed to publish mapping event to delete permanently. But the project(%s) was deleted successfully. Err - %v", req.GetLuid(), err)
	}

	return &projectpb.DeleteProjectResponse{Luid: req.GetLuid(), DeletedAt: pkgUtils.DbFormatDate(time.Now())}, nil
}

func (p *ProjectServer) GetProjectVersions(_ context.Context, req *projectpb.ProjectLuidRequest) (*projectpb.ProjectVersionsResponse, error) {
	versions := make([]*projectpb.Version, 0)
	_, err := p.DB.Query(&versions, fmt.Sprintf(constants.QueryFetchProjectVersionsByLuid, req.GetLuid()))

	if err == pg.ErrNoRows {
		return nil, utils.ProjectNotFoundError(req.GetLuid())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch project with id : %s from db while getting project-versions, internal error, %v", req.GetLuid(), err))
	}

	return &projectpb.ProjectVersionsResponse{Versions: versions}, nil
}

func (p *ProjectServer) DeleteComponent(ctx context.Context, req *projectpb.DeleteComponentRequest) (*projectpb.DeleteComponentResponse, error) {
	// get project from cache
	project, err := p.GetProjectFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if not in cache, fetch from DB
	if len(project.GetLuid()) == 0 {
		if project, err = p.GetCurrentActiveProject(req.GetProjectId()); err != nil {
			return nil, err
		}
	}

	// if the project not found in db or cache, return Not Found
	if len(project.GetLuid()) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Project not found in cache or db with project id(%s)", req.GetProjectId()))
	}

	// if the component not found in project key in cache, return Not Found
	if project.GetConfig().GetComponents()[req.GetComponentId()] == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Component not found with id(%s) in project with id in(%s)", req.GetComponentId(), req.GetProjectId()))
	}

	dependencies := make(map[string]string)
	// Check the component dependencies
	for compId, comp := range project.GetConfig().GetComponents() {
		if compId != req.GetComponentId() {
			compStr, err := utils.MarshalToString(comp)
			if err != nil {
				return nil, status.Error(codes.Unknown, fmt.Sprintf("Failed in string conversion while checking for the dependencies for the component id(%s)", compId))
			}
			if strings.Contains(compStr, req.GetComponentId()) {
				dependencies[compId] = ""
			}
		}
	}

	// update the dependency names and send the dependencies
	if len(dependencies) > 0 {
		utils.UpdateCompNames(dependencies, utils.ConvertComponentMapToArr(project.GetConfig().GetComponents()))
		return &projectpb.DeleteComponentResponse{
			ProjectId:             req.GetProjectId(),
			ComponentDependencies: utils.MapToArrayDependencies(dependencies),
			Components:            utils.ConvertComponentMapToArr(project.GetConfig().GetComponents()),
		}, nil
	}

	//update the cache with changes
	deleted := false
	project.Config.Components = utils.ConvertComponentArrToMap(utils.DeleteComponentFromComponentArr(req.GetComponentId(), utils.ConvertComponentMapToArr(project.GetConfig().GetComponents()), &deleted))
	delete(project.Config.Components, req.GetComponentId())

	err = p.SetProjectInCache(ctx, project)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete component(%s) in project(%s) from cache. Error - %v", req.GetComponentId(), req.GetProjectId(), err))
	}

	return &projectpb.DeleteComponentResponse{
		ProjectId:  req.GetProjectId(),
		UpdatedAt:  pkgUtils.DbFormatDate(time.Now()),
		Components: utils.ConvertComponentMapToArr(project.GetConfig().GetComponents()),
	}, nil
}
