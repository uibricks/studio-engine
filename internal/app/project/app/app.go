package app

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	project "github.com/uibricks/studio-engine/internal/app/project/services"
	"github.com/uibricks/studio-engine/internal/pkg/config"
	"github.com/uibricks/studio-engine/internal/pkg/db"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/middleware"
	"github.com/uibricks/studio-engine/internal/pkg/proto/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	ServerConfig     config.ServerConfig
	DbConfig         config.DBConfig
	RedisConfig      config.RedisConfig
	RabbitMqConfig   config.RabbitMqConfig
	PrometheusConfig config.PrometheusConfig
	DbClient         db.DbClient
	GrpcServerOpts   []grpc.ServerOption
	ProjectServer    *project.ProjectServer
}

func (a *App) Start(check func(error)) {
	lis, err := net.Listen("tcp", a.ServerConfig.ServiceAddress)
	check(err)

	err = a.DbClient.Connect(a.DbConfig.Url, a.DbConfig.Schema)
	check(err)

	s := grpc.NewServer(a.GrpcServerOpts...)

	service.RegisterProjectServiceServer(s, a.ProjectServer)

	grpc_prometheus.Register(s)
	middleware.RunPrometheusServer(a.PrometheusConfig)

	logger.Log.Info("Starting project service...")
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Log.Fatal("Failed to start the service", zap.Error(err))
		}
	}()
}

func (a *App) Shutdown() {
	a.DbClient.Close()
}
