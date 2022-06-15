package constants

import (
	"github.com/teris-io/shortid"
	"github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"time"
)

const (
	DefaultQueueExpiration time.Duration = 300
)

const (
	DefaultAuthenticationType string = "noauth"
	BodyTypeForm string = "form"
	BodyTypeJson string = "json"
	NotFound string = "NotFound"
)

const (
	ContentTypeJson string = "application/json"
	HeaderKeyContentType string = "content-type"
	QueryParamsSeparator string = "?"
	QueryParamAppend string = "&"
	QueryParamValueSeparator string = "="
	EnvVarPrefix string = "[env."
	EnvVarSuffix string = "]"
	UrlSuffix string = "/"
	GroupExpressionType = "group"
)

func GetReplyQueueName() constants.ReplyQueueName {
	sid, err := shortid.Generate()
	if err != nil {
		logger.Sugar.Errorf("Failed to generate the short id for reply message queue name")
	}
	return constants.ReplyQueueName("mapping_reply" + sid)
}