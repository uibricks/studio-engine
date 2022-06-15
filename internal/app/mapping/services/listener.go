package mapping

import (
	"context"
	"encoding/json"
	"github.com/uibricks/studio-engine/internal/app/mapping/constants"
	pkgConstants "github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	globalutils "github.com/uibricks/studio-engine/internal/pkg/utils"
)

// Listen to mapping queue and handles message on arrival
// Activated on service startup and runs indefinitely in a go routine
func (m *MappingServer) Listen() {
	go m.messageHandler(context.Background())
}

// processes messages received on mapping-queue - runs on a infinite for loop in a separate go-routine
// acknowledges request that match action
// rejects messages otherwise
func (m *MappingServer) messageHandler(ctx context.Context) {

	if _, err := rabbitmq.ProvideQueueWithExp(pkgConstants.ReplyQueueName(rabbitmq.Mapping_Queue_Name), constants.DefaultQueueExpiration, m.Channel); err != nil {
		logger.Sugar.Errorf("Failed to declare mapping queue - %v", err)
	}

	for {
		// continuously consumes from Mapping Queue
		d, err := m.Rabbit.Consume(rabbitmq.Mapping_Queue_Name, false, m.Channel)
		if err != nil {
			logger.Sugar.Errorf("Failed to fetch messages from mapping queue - %v", err)
		}

		for msg := range d.Messages {

			rmqMsg := &rabbitmq.RMQMessage{}
			json.Unmarshal(msg.Body, rmqMsg)

			context := globalutils.ReconstructContext(rmqMsg.SharedContext)

			err := m.ExecuteAction(context, rmqMsg)

			msg.Ack(false)

			// reply only if a callback queue is available
			if msg.ReplyTo != "" {
				status := rabbitmq.Status_Error
				if err == nil {
					status = rabbitmq.Status_Success
				} else {
					logger.WithContext(context).Errorf("failed to %s - %v", rmqMsg.Action, err)
				}

				b, _ := m.Rabbit.PrepareRMQMessage(nil, rmqMsg.Action, status)

				err = m.Rabbit.Publish(msg.ReplyTo, msg.CorrelationId, b, m.Channel)

				if err != nil {
					logger.WithContext(context).Errorf("failed to publish message to queue for %s - %v", rmqMsg.Action, err)
				}
			}
		}
	}
}
