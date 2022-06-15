package expression

import (
	"context"
	"encoding/json"
	"fmt"
	expressionpb "github.com/uibricks/studio-engine/internal/pkg/proto/expression"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ExecuteAction(ctx context.Context, rmqMsg *rabbitmq.RMQMessage) (string, error) {

	switch rmqMsg.Action {

	case rabbitmq.Action_Resolve_Expression:

		req := &expressionpb.EvalExpressionRequest{}
		b, err := json.Marshal(rmqMsg.Payload)
		if err != nil {
			return "", status.Error(codes.Internal, fmt.Sprintf("failed to marshal payload to string payload(%v),%v", rmqMsg.Payload, err))
		}

		err = json.Unmarshal(b, req)
		if err != nil {
			return "", fmt.Errorf("error while invoking action : %s, %v", rmqMsg.Action, err)
		}

		resp, err := s.EvalExpression(ctx, req)

		if err != nil {
			return "", err
		}

		return resp.ExpressionResult, err
	default:
		return "", fmt.Errorf("invalid action sent in expression queue : %s", rmqMsg.Action)
	}

}
