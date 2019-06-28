package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	v1pb "github.com/vterdunov/janna-proto/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/config"
	deliveryGrpc "github.com/vterdunov/janna/internal/delivery/grpc"
	"github.com/vterdunov/janna/internal/delivery/grpc/middleware"
	vmWareRepository "github.com/vterdunov/janna/internal/virtualmachine/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Could not read config. Err: %s\n", err)
		os.Exit(1)
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// setup GRPC server with middlewares
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			middleware.LoggingInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	// setup repositories
	appRep := appinfo.NewAppRepository()
	vmwareRep, err := vmWareRepository.NewVMRepository(cfg.VMWare.URL, cfg.VMWare.Insecure)
	if err != nil {
		logger.Error().Err(err).Msg("could not create VMWare connection")
		os.Exit(1)
	}

	// register and run servers
	deliveryGrpc.RegisterServer(grpcServer, appRep, vmwareRep)

	var httpServer *http.Server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run HTTP server
	logger.Info().Msg("starting HTTP Gateway proxy...")
	httpServer = setupHTTPServer(ctx, cfg, &logger)
	go func() {
		httpServer.ListenAndServe()
	}()

	// run GRPC server
	logger.Info().Msg("starting GRPC server...")
	l, err := net.Listen("tcp", ":"+cfg.Protocols.GRPC.Port) //nolint:gosec
	if err != nil {
		logger.Error().Err(err).Msg("could not start GRPC server")
	}

	go func() {
		if err = grpcServer.Serve(l); err != nil {
			logger.Error().Err(err).Msg("unenxpected error")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	// graceful shutdown
	logger.Info().Msg("shutting down HTTP Gateway proxy...")
	httpServer.Shutdown(ctx)
	logger.Info().Msg("shutting down gRPC server...")
	grpcServer.GracefulStop()

}

func setupHTTPServer(ctx context.Context, cfg *config.Config, l *zerolog.Logger) *http.Server {
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(populateXRequestID),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	grpcPort := cfg.Protocols.GRPC.Port
	if err := v1pb.RegisterJannaAPIHandlerFromEndpoint(ctx, gwMux, ":"+grpcPort, opts); err != nil {
		l.Error().Err(err).Msg("failed to start HTTP gateway")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", contextWrap(gwMux))

	server := http.Server{
		Addr:    ":" + cfg.Protocols.HTTP.Port,
		Handler: mux,
	}

	return &server
}

func contextWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestID := req.Header.Get("X-Request-Id")
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		h.ServeHTTP(w, req.WithContext(ctx))
	})
}

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

func populateXRequestID(ctx context.Context, req *http.Request) metadata.MD {
	m := map[string]string{}
	reqID, ok := ctx.Value(RequestIDKey).(string)
	if ok && reqID != "" {
		m["request_id"] = reqID
		return metadata.New(m)
	}

	id := uuid.New()
	m["request_id"] = id.String()

	return metadata.New(m)
}
