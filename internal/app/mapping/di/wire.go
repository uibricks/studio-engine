//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/uibricks/studio-engine/internal/app/mapping/app"
	"github.com/uibricks/studio-engine/internal/app/mapping/config"
	mapping "github.com/uibricks/studio-engine/internal/app/mapping/services"
	configPkg "github.com/uibricks/studio-engine/internal/pkg/config"
	"github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/db"
	"github.com/uibricks/studio-engine/internal/pkg/middleware"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"github.com/uibricks/studio-engine/internal/pkg/redis"
)

var configSet = wire.NewSet(
	config.ProvideAppConfig,
	wire.FieldsOf(new(config.AppConfig), "ServerConfig"),
	wire.FieldsOf(new(config.AppConfig), "DatabaseConfig"),
	wire.FieldsOf(new(config.AppConfig), "RedisConfig"),
	wire.FieldsOf(new(config.AppConfig), "RabbitMqConfig"),
	wire.FieldsOf(new(configPkg.RabbitMqConfig), "QueuePrefix"),
	wire.FieldsOf(new(config.AppConfig), "PrometheusConfig"),
)

var rabbitMqSet = wire.NewSet(
	rabbitmq.ProvideDefaultRabbitMqConn,
	rabbitmq.ProvideQueue,
	rabbitmq.ProvideChannel,
	rabbitmq.ProvideQueueWithExp,
	constants.ProvideReplyQueueName,
)

func InitializeApp() (*app.App, error) {
	wire.Build(
		configSet,
		rabbitMqSet,
		db.ProvideDBClient,
		middleware.ProvideGrpcServerOpts,
		mapping.ProvideMappingServer,
		redis.ProvideDefaultRedisConn,
		wire.Struct(new(app.App), "*"),
	)

	return &app.App{}, nil
}
