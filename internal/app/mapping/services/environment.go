package mapping

import (
	"context"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	globalutils "github.com/uibricks/studio-engine/internal/pkg/utils"
	"time"
)

func (m *MappingServer) UpdateEnvironment(ctx context.Context, req *mappingpb.UpdateEnvRequest) (*mappingpb.UpdateEnvResponse, error) {

	// get mapping from cache
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if not in cache, fetch from DB : both projectId and projectVersion are required
	if len(mapping.GetProjectId()) == 0 {
		if err = m.LoadMappingWithVersion(mapping, req.GetProjectId(), req.GetProjectVersion()); err != nil {
			return nil, err
		}
	}

	// update mapping object with data from request
	if len(req.DefaultEnvironment) > 0 {
		mapping.Config.DefaultEnvironment = req.GetDefaultEnvironment()
	}
	if req.Environments != nil || req.GetEmptyEnvironments() {
		mapping.Config.Environments = req.GetEnvironments()
	}
	if req.EnvironmentVariables != nil || req.GetEmptyEnvironmentVariables() {
		mapping.Config.EnvironmentVariables = req.GetEnvironmentVariables()
	}

	// set default environment to empty if there are no environments
	if mapping.GetConfig().GetEnvironments() == nil {
		mapping.Config.DefaultEnvironment = ""
	}

	// set updated mapping back to cache.
	err = m.SetMappingInCache(ctx, mapping, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	return &mappingpb.UpdateEnvResponse{
		ProjectId:            req.GetProjectId(),
		UpdatedAt:            globalutils.DbFormatDate(time.Now()),
		DefaultEnvironment:   mapping.GetConfig().GetDefaultEnvironment(),
		Environments:         mapping.GetConfig().GetEnvironments(),
		EnvironmentVariables: mapping.GetConfig().GetEnvironmentVariables(),
	}, nil
}
