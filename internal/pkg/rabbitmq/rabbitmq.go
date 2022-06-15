package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/uibricks/studio-engine/internal/pkg/config"
	constantsPkg "github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/env"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	"github.com/uibricks/studio-engine/internal/pkg/types"
	"github.com/uibricks/studio-engine/internal/pkg/utils"
	"go.uber.org/zap"
	"time"
)

type Rabbit struct {
	conn *amqp.Connection
}

// Channel amqp.ProvideChannel wapper
type Channel struct {
	*amqp.Channel
	closed bool
}

type Params struct {
	rabbitHost string
	rabbitPort string
	rabbitUser string
	rabbitPwd  string
	rabbitEncrypted bool
}

var connParams *Params = &Params{}

type RMQMessage struct {
	Action        string
	Payload       interface{}
	SharedContext types.SharedContext
}

type ConsumeDelivery struct {
	Messages <-chan amqp.Delivery
	Channel  *Channel
}

func ProvideDefaultRabbitMqConn(config config.RabbitMqConfig) (*Rabbit, error) {
	connParams.rabbitHost = config.Host
	connParams.rabbitPort = config.Port
	connParams.rabbitUser = config.User
	connParams.rabbitPwd = config.Password
	connParams.rabbitEncrypted = env.GetBool("rabbit_encrypted")
	return dial()
}

// Channel wrap amqp.Connection.ProvideChannel, get a auto reconnect channel
func ProvideChannel(r *Rabbit) (*Channel, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	ch.Qos(1, 0, false)

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-ch.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.closed {
				channel.Close()
				break
			}

			logger.Log.Debug("channel closed", zap.Any("reason", reason))

			for {
				time.Sleep(reconnect_Delay_In_Seconds * time.Second)

				ch, err := r.conn.Channel()
				if err == nil {
					channel.Channel = ch
					logger.Log.Debug("channel successfully reconnected")
					break
				}

				logger.Log.Debug("channel recreation failed", zap.Error(err))

			}
		}
	}()

	return channel, nil
}

// ProvideQueueWithExp - idempotent. This will close the channel after declaring a queue and it will
// set the expiration time for a message
type QueueWithExpiry amqp.Queue
func ProvideQueueWithExp(name constantsPkg.ReplyQueueName, secs time.Duration, ch *Channel) (QueueWithExpiry, error) {
	args := make(amqp.Table)
	args["x-message-ttl"] = int32(secs * 1000)

	q, err := declareQueue(name, args, ch)
	if err != nil {
		return QueueWithExpiry(amqp.Queue{}), err
	}
	return QueueWithExpiry(q), nil
}

// declare queue - idempotent. This will close the channel after declaring a queue
func ProvideQueue(name constantsPkg.ReplyQueueName, ch *Channel) (amqp.Queue, error) {
	return declareQueue(name, nil, ch)
}

func Conn(host string, port string, user string, pwd string, encrypted bool) (*Rabbit, error) {
	connParams.rabbitHost = host
	connParams.rabbitPort = port
	connParams.rabbitUser = user
	connParams.rabbitPwd = pwd
	connParams.rabbitEncrypted = encrypted
	return dial()
}

func (c *Channel) Close() error {
	if c.closed {
		return amqp.ErrClosed
	}

	c.closed = true
	return c.Channel.Close()
}

func url() string {
	protocol := "amqp"
	if connParams.rabbitEncrypted{
		protocol += "s"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/",protocol, connParams.rabbitUser, connParams.rabbitPwd, connParams.rabbitHost, connParams.rabbitPort)
}

func dial() (*Rabbit, error) {
	conn, err := amqp.Dial(url())

	if err != nil {
		return nil, &ConnectionError{err: err}
	}

	rabbit := &Rabbit{conn: conn}

	go func() {
		for {
			reason, ok := <-rabbit.conn.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				logger.Sugar.Debug("Rabbitmq connection closed")
				break
			}

			logger.Sugar.Debugf("Rabbitmq connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait for reconect
				time.Sleep(reconnect_Delay_In_Seconds * time.Second)

				conn, err := amqp.Dial(url())
				if err == nil {
					rabbit.conn = conn
					logger.Sugar.Debug("Rabbitmq connection reestablished successfully")
					break
				}
				logger.Sugar.Debugf("Reconnecting Rabbitmq server failed, err:%v", err)
			}
		}
	}()

	return rabbit, nil
}

// declare queue - idempotent. This will close the channel after declaring a queue
func declareQueue(name constantsPkg.ReplyQueueName, args amqp.Table, ch *Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(string(name),
		true, false, false, false, args)

	if err != nil {
		return q, &DeclareQueueError{name: string(name), err: err}
	}
	return q, nil
}

func (r *Rabbit) Consume(queueName string, autoAck bool, ch *Channel) (*ConsumeDelivery, error) {
	msgs, err := ch.Consume(queueName, "", autoAck, false, false, false, nil)

	if err != nil {
		return nil, &ConsumeError{name: queueName, err: err}
	}

	d := &ConsumeDelivery{
		Messages: msgs,
		Channel:  ch,
	}

	return d, nil
}

// ConsumeRPCMessage - This will consume a message from queue and compares the correlation id and action with the message. If everything
// matches then it will return the message and ackowledges the queue. If it doens't match then it will keep continues to receive the message till timeout
func (r *Rabbit) ConsumeRPCMessage(ctx context.Context, queueName string, action string, corrId string, ch *Channel) (*RMQMessage, error) {
	d, err := r.Consume(queueName, false, ch)
	defer d.Channel.Close()

	if err != nil {
		return nil, err
	}

	select {
	case m := <-d.Messages:

		rmqMsg := &RMQMessage{}
		json.Unmarshal(m.Body, rmqMsg)

		if rmqMsg.Action == action {
			if corrId == m.CorrelationId {
				if err := m.Ack(false); err != nil {
					logger.WithContext(ctx).Errorf("failed to acknowledge queue(%s), corrId(%s), action(%s) - %v", queueName, corrId, action, err)
				}
				return rmqMsg, nil
			}
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout - failed to consume message queue(%s), corrId(%s), action(%s)", queueName, corrId, action)
	}

	return nil, nil
}

func (r *Rabbit) PublishWithCallBack(queueName string, callbackQueue string, corrId string, res []byte, ch *Channel) error {
	err := ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			ReplyTo:       callbackQueue,
			CorrelationId: corrId,
			Body:          res,
		})

	return err
}

func (r *Rabbit) Publish(queueName string, corrId string, res []byte,ch *Channel) error {
	err := ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			CorrelationId: corrId,
			Body:          res,
		})

	return err
}

func (r *Rabbit) PrepareRMQMessage(ctx context.Context, action string, payload interface{}) ([]byte, error) {
	sharedContext := types.SharedContext{}
	if ctx != nil {
		sharedContext.RequestID = request.GetContextRequestID(ctx)
		ctxData, err := utils.PrepareSharableContext(&ctx, utils.SerializeOpts{RetainCancel: true, RetainDeadline: true})
		if err == nil {
			sharedContext.ContextData = ctxData
		}
	}
	b, err := json.Marshal(&RMQMessage{Action: action, Payload: payload, SharedContext: sharedContext})
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to marshal payload with action (%s) and payload (%v)", action, payload)
	}
	return b, err
}
