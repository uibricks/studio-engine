package middleware

import (
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"google.golang.org/grpc"
)

func AddInterceptors(opts []grpc.ServerOption, uInterceptors []grpc.UnaryServerInterceptor, sInterceptors []grpc.StreamServerInterceptor) []grpc.ServerOption {
	opts = append(opts, grpc.ChainUnaryInterceptor(uInterceptors...))
	opts = append(opts, grpc.ChainStreamInterceptor(sInterceptors...))
	return opts
}

// Add grpc default middleware like logging and prometheus metrics
func ProvideGrpcServerOpts() []grpc.ServerOption {
	// gRPC server startup options
	opts := []grpc.ServerOption{}

	uInterceptors := []grpc.UnaryServerInterceptor{}
	sInterceptors := []grpc.StreamServerInterceptor{}

	// add middleware
	AddLogging(logger.Log, &uInterceptors, &sInterceptors)
	AddPrometheus(&uInterceptors, &sInterceptors)
	AddRecovery(&uInterceptors, &sInterceptors)

	opts = AddInterceptors(opts, uInterceptors, sInterceptors)

	return opts
}
