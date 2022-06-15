package expression

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uibricks/studio-engine/internal/app/expression/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	globalutils "github.com/uibricks/studio-engine/internal/pkg/utils"
)

// Listen to expression queue and handles message on arrival
// Activated on service startup and runs indefinitely in a go routine
func (s *Server) Listen() {
	go s.messageHandler(context.Background())
}

// processes messages received on expression-queue - runs on a infinite for loop in a separate go-routine
// acknowledges request that match action
// rejects messages otherwise
func (s *Server) messageHandler(ctx context.Context) {

	if _, err := rabbitmq.ProvideQueueWithExp(rabbitmq.Expression_Queue_Name, constants.DefaultQueueExpiration, s.Channel); err != nil {
		logger.Sugar.Errorf("Failed to declare expression queue - %v", err)
	}

	for {
		// continuously consumes from Expression Queue
		d, err := s.Rabbit.Consume(rabbitmq.Expression_Queue_Name, false, s.Channel)
		if err != nil {
			logger.Sugar.Errorf("Failed to fetch messages from expression queue - %v", err)
		}

		for msg := range d.Messages {

			rmqMsg := &rabbitmq.RMQMessage{}
			err := json.Unmarshal(msg.Body, rmqMsg)
			if err != nil {
				logger.WithContext(context.Background()).Error(fmt.Sprintf("Error while invoking action : %s, %v", rmqMsg.Action, err))
			}

			ctx := globalutils.ReconstructContext(rmqMsg.SharedContext)

			resp, err := s.ExecuteAction(ctx, rmqMsg)

			err = msg.Ack(false)
			if err != nil {
				logger.WithContext(ctx).Error(fmt.Sprintf("Error while acknowledging expression rmq message for action : %s, %v", rmqMsg.Action, err))
			}

			// reply only if a callback queue is available
			if msg.ReplyTo != "" {
				var status string
				if err == nil {
					status = resp
				} else {
					status = fmt.Sprintf("failed to %s - %v", rmqMsg.Action, err)
					logger.WithContext(ctx).Error(status)
				}

				b, _ := s.Rabbit.PrepareRMQMessage(nil, rmqMsg.Action, status)

				err = s.Rabbit.Publish( msg.ReplyTo, msg.CorrelationId, b, s.Channel)

				if err != nil {
					logger.WithContext(ctx).Errorf("failed to publish message to queue for %s - %v", rmqMsg.Action, err)
				}
			}
		}
	}
}
