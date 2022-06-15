package middleware

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uibricks/studio-engine/internal/pkg/config"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

//var (
//	// Creates a metrics registry
//	reg = prometheus.NewRegistry()
//
//	// Create some standard server metrics.
//	grpcMetrics = grpc_prometheus.NewServerMetrics()
//)
//
//func CreateCounterMetric(name string, help string) *prometheus.CounterVec {
//	counterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
//		Name: name,
//		Help: help,
//	}, []string{"name"})
//
//	//Register standard server metrics and customized metrics to registry
//	reg.MustRegister(grpcMetrics, counterMetric)
//
//	return counterMetric
//}

func AddPrometheus(uInterceptors *[]grpc.UnaryServerInterceptor, sInterceptors *[]grpc.StreamServerInterceptor) {
	*uInterceptors = append(*uInterceptors, grpc_prometheus.UnaryServerInterceptor)
	*sInterceptors = append(*sInterceptors, grpc_prometheus.StreamServerInterceptor)
}

//func InitializeAllMetrics(server *grpc.Server) {
//	grpcMetrics.InitializeMetrics(server)
//}

func registerPrometheus(server *grpc.Server) {
	grpc_prometheus.Register(server)
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func RunPrometheusServer(config config.PrometheusConfig) {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
			log.Fatalf("Unable to start a http server for prometheus. %v", err)
		}
	}()
}
