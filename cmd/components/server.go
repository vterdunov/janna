package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1pb "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/config"
	deliveryGrpc "github.com/vterdunov/janna/internal/delivery/grpc"
	"github.com/vterdunov/janna/internal/delivery/grpc/middleware"
	"github.com/vterdunov/janna/internal/virtualmachine/repository"
	"github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Could not read config. Err: %s\n", err)
		os.Exit(1)
	}

	// Create logger, which we'll use and give to other components.
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Could not create logger. Err: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	grpc_zap.ReplaceGrpcLogger(logger)

	// setup GRPC server with middlewares
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			// grpc_recovery.UnaryServerInterceptor(),
			middleware.NoopInterceptor,
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	// setup repositories
	appRep := repository.NewAppRepository()
	vmwareRep, err := repository.NewVMWareRepository(cfg.VMWare.URL, cfg.VMWare.Insecure)
	if err != nil {
		logger.Error("could not create VMWare connection", zap.Error(err))
		os.Exit(1)
	}

	// register and run servers
	// usecase := usecase.NewUsecase(appRep, vmwareRep)
	appInfoUsecase := usecase.AppInfo{appRep}
	vmInfoUsecase := usecase.VMInfo{vmwareRep}
	vmDeployUsecase := usecase.VMDeploy{vmwareRep}
	deliveryGrpc.RegisterServer(grpcServer, &appInfoUsecase, &vmInfoUsecase, &vmDeployUsecase)

	var httpServer *http.Server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run HTTP server
	logger.Info("starting HTTP Gateway proxy...")
	httpServer = setupHTTPServer(ctx, cfg, logger)
	go func() {
		httpServer.ListenAndServe()
	}()

	logger.Info("starting GRPC server...")
	l, err := net.Listen("tcp", ":"+cfg.Protocols.GRPC.Port) //nolint:gosec
	if err != nil {
		logger.Error("could not start GRPC server", zap.Error(err))
	}

	// run GRPC server
	go func() {
		if err = grpcServer.Serve(l); err != nil {
			logger.Error("unenxpected error", zap.Error(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	// graceful shutdown
	logger.Info("shutting down HTTP Gateway proxy...")
	httpServer.Shutdown(ctx)
	logger.Info("shutting down gRPC server...")
	grpcServer.GracefulStop()

}

func setupHTTPServer(ctx context.Context, cfg *config.Config, l *zap.Logger) *http.Server {
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	grpcPort := cfg.Protocols.GRPC.Port
	if err := v1pb.RegisterJannaAPIHandlerFromEndpoint(ctx, gwMux, ":"+grpcPort, opts); err != nil {
		l.Error("failed to start HTTP gateway", zap.Error(err))
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", gwMux)

	server := http.Server{
		Addr:    ":" + cfg.Protocols.HTTP.Port,
		Handler: mux,
	}

	return &server
}
