package utils

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ProjectNotFoundError(projectID string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("Project not found with id - %s", projectID))
}
