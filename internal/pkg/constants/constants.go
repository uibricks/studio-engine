package constants

import (
	"fmt"
	"github.com/teris-io/shortid"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
)

type RedisDatabase int

const (
	Project = iota
	Mapping
)

const (
	ProjectTypeRecent string = "recent"
	ProjectTypeTrash         = "trash"
	ProjectTypeAll           = "all"
)

const (
	Saga_Timeout_In_Minutes = 5
)

const (
	STATE_ACTIVE   string = "active"
	STATE_INACTIVE        = "inactive"
	STATE_PENDING         = "pending"
	STATE_PANIC           = "panic"
	STATE_CACHE           = "cache"
)

type ReplyQueueName string

func ProvideReplyQueueName(queuePrefix string) ReplyQueueName {
	sid, err := shortid.Generate()
	if err != nil {
		logger.Sugar.Errorf("Failed to generate the short id for reply message queue name")
	}
	return ReplyQueueName(fmt.Sprintf("%s_%s", queuePrefix, sid))
}
