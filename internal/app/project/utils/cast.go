package utils

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	projectPb "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CastStringToProject(val string) (*projectPb.Object, error) {
	project := &projectPb.Object{}
	if len(val) > 0 {
		err := jsonpb.UnmarshalString(val, project)
		if err != nil {
			return nil, err
		}
	}
	return project, nil
}

func CastStringToProjects(val interface{}) ([]*projectPb.Object, error) {
	obj := val.([]interface{})
	projects := []*projectPb.Object{}

	for i := 1; i <= len(obj); i += 2 {
		project := &projectPb.Object{}
		err := jsonpb.UnmarshalString(obj[i].(string), project)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to unmarshall the project, Error - %v", err))
		}
		projects = append(projects, project)
	}
	return projects, nil

}
