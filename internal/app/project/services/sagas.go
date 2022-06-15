package project

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/uibricks/studio-engine/internal/app/project/constants"
	pkgConstants "github.com/uibricks/studio-engine/internal/pkg/constants"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	projectpb "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	"github.com/uibricks/studio-engine/internal/pkg/saga"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *ProjectServer) saveProjectTransaction(ctx context.Context, req *projectpb.SaveProjectRequest, project *projectpb.Object) error {
	project.UserVersion = req.GetUserVersion()
	project.State = pkgConstants.STATE_PENDING
	_, err := p.DB.Model(project).Insert()
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("Failed to save the project, internal error, %v", err))
	}
	return nil
}

func (p *ProjectServer) saveMappingTransaction(ctx context.Context, req *projectpb.SaveProjectRequest, project *projectpb.Object) error {
	corrId := request.GetContextRequestID(ctx)

	// get project current version
	currPrj, err := p.GetCurrentActiveProject(project.GetLuid())

	// if there are no rows then dont throw the error
	if err != pg.ErrNoRows && err != nil {
		return fmt.Errorf("failed to read the project current version while save mapping transaction in save project saga, error - %v", err)
	}

	// prepare payload
	b, _ := p.Rabbit.PrepareRMQMessage(ctx,
		rabbitmq.Action_Save_Mapping, &mappingpb.SaveMappingRequest{
			ProjectId:          req.GetLuid(),
			NewProjectVersion:  project.GetVersion(),
			CurrProjectVersion: currPrj.GetVersion(),
		})

	// publish to mapping queue
	err = p.Rabbit.PublishWithCallBack(ctx, rabbitmq.Mapping_Queue_Name, p.ReplyQueue.Name, corrId, b)
	if err != nil {
		return fmt.Errorf("failed to publish message to mapping queue - %v", err)
	}

	// consume message
	msg, err := p.Rabbit.ConsumeRPCMessage(ctx, p.ReplyQueue.Name, rabbitmq.Action_Save_Mapping, corrId)
	if err != nil {
		return err
	}

	status := fmt.Sprintf("%s", msg.Payload)
	if status == rabbitmq.Status_Error {
		return fmt.Errorf("failed to save the mapping to db")
	}

	return nil
}

func (p *ProjectServer) rollbackSaveProject(ctx context.Context, _ *projectpb.SaveProjectRequest, project *projectpb.Object) error {
	if err := p.UpdateProjectStatus(ctx, pkgConstants.STATE_INACTIVE, project); err != nil {
		return err
	}
	return nil
}

func (p *ProjectServer) SaveProjectSaga(ctx context.Context, req *projectpb.SaveProjectRequest, project *projectpb.Object) error {

	sec := saga.AddSubTxDef(constants.TxSaveProject, p.saveProjectTransaction, nil)
	sec.AddSubTxDef(constants.TxSaveMapping, p.saveMappingTransaction, p.rollbackSaveProject)
	sec.AddSubTxDef(constants.TxCommitProject, p.UpdateProjectStatus, nil)

	sg := sec.StartSaga(ctx)

	sg, res := sg.ExecSub(constants.TxSaveProject, ctx, req, project)
	if err := saga.GetError(res); err != nil {
		sg.EndSaga()
		return err
	}

	sg, res = sg.ExecSub(constants.TxSaveMapping, ctx, req, project)
	if err := saga.GetError(res); err != nil {
		sg.EndSaga()
		return err
	}

	sg, res = sg.ExecSub(constants.TxCommitProject, ctx, pkgConstants.STATE_ACTIVE, project)
	if err := saga.GetError(res); err != nil {
		sg.EndSaga()
		return err
	}

	// Triggers log cleanup, else on rollback all previous saga's rollbacks are also getting executed
	sg.EndSaga()

	return nil
}

func (p *ProjectServer) restoreProjectTransaction(_ context.Context, req *projectpb.ProjectLuidRequest) error {
	rows, err := p.DB.Model(&projectpb.Object{}).Set("deleted_at=null").Where(fmt.Sprintf("luid::text='%s'", req.GetLuid())).Update()
	if err != nil || rows.RowsAffected() == 0 {
		return status.Error(codes.Internal, fmt.Sprintf("Failed to restore project, %v", err))
	}
	return nil
}

func (p *ProjectServer) restoreMappingTransaction(ctx context.Context, req *projectpb.ProjectLuidRequest) error {
	corrId := request.GetContextRequestID(ctx)

	b, _ := p.Rabbit.PrepareRMQMessage(ctx,
		rabbitmq.Action_Restore_Mapping, &mappingpb.RestoreMappingRequest{
			ProjectId:      req.GetLuid(),
		})

	// publish to mapping queue
	err := p.Rabbit.PublishWithCallBack(ctx, rabbitmq.Mapping_Queue_Name, p.ReplyQueue.Name, corrId, b)
	if err != nil {
		return fmt.Errorf("failed to publish message to mapping queue - %v", err)
	}

	// consume message
	msg, err := p.Rabbit.ConsumeRPCMessage(ctx, p.ReplyQueue.Name, rabbitmq.Action_Restore_Mapping, corrId)
	if err != nil {
		return err
	}

	status := fmt.Sprintf("%s", msg.Payload)
	if status == rabbitmq.Status_Error {
		return fmt.Errorf("failed to restore the mapping in db")
	}
	return nil
}

func (p *ProjectServer) RestoreProjectSaga(ctx context.Context, req *projectpb.ProjectLuidRequest) error {

	sec := saga.AddSubTxDef(constants.TxRestoreMapping, p.restoreMappingTransaction, nil)
	// even if mapping is restored and project restoration fails, api-gateway will get an error response
	// since project is the driving point mapping will practically remain deleted
	sec.AddSubTxDef(constants.TxRestoreProject, p.restoreProjectTransaction, nil)

	sg := sec.StartSaga(ctx)

	sg, res := sg.ExecSub(constants.TxRestoreMapping, ctx, req)
	if err := saga.GetError(res); err != nil {
		sg.EndSaga()
		return err
	}

	sg, res = sg.ExecSub(constants.TxRestoreProject, ctx, req)
	if err := saga.GetError(res); err != nil {
		sg.EndSaga()
		return err
	}

	// Triggers log cleanup, else on rollback all previous saga's rollbacks are also getting executed
	sg.EndSaga()

	return nil
}
