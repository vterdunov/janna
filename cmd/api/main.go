package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	v1pb "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/config"
	deliveryGrpc "github.com/vterdunov/janna/internal/delivery/grpc"
	"github.com/vterdunov/janna/internal/delivery/grpc/middleware"
	"github.com/vterdunov/janna/internal/log"
	"github.com/vterdunov/janna/internal/producer"
)

func main() {
	logger := log.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Error(err, "could not read config")
		os.Exit(1)
	}

	// setup GRPC server with middlewares
	var grpcMiddlewares []grpc.UnaryServerInterceptor
	grpcMiddlewares = append(grpcMiddlewares, grpc_prometheus.UnaryServerInterceptor)
	if !cfg.Debug {
		grpcMiddlewares = append(grpcMiddlewares, grpc_recovery.UnaryServerInterceptor())
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcMiddlewares...,
		)),
	)

	// register metrics
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	// create publisher
	producer, err := producer.NewProducer("redis://redis:6379")
	if err != nil {
		logger.Error(err, "could not create Worker")
		os.Exit(1)
	}

	// create TCP listener
	ln, err := net.Listen("tcp", ":"+cfg.Protocols.HTTP.Port)
	if err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// setup repositories
	appRep := appinfo.NewAppRepository()

	// register service and middlewares
	service := deliveryGrpc.NewService(appRep, producer)
	service = middleware.NewLoggingMiddleware(service, logger)
	deliveryGrpc.RegisterServer(grpcServer, service, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := createHTTPServer(ctx, ln, grpcServer)
	if err != nil {
		logger.Error(err, "could not create http server")
		os.Exit(1)
	}

	l := logger.WithFields("addr", ln.Addr().String())
	l.Info("start listening address")

	fin := make(chan struct{})
	go func() {
		<-ctx.Done()
		logger.Info("shutting down shared server...")
		_ = srv.Shutdown(ctx)
		close(fin)
	}()

	go func() {
		err = srv.Serve(ln)
		if err != nil {
			logger.Error(err, "ListenAndServe")
			os.Exit(1)
		}
	}()

	<-fin
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if r.ProtoMajor == 2 && strings.Contains(contentType, "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func createHTTPServer(ctx context.Context, ln net.Listener, grpcServer *grpc.Server) (*http.Server, error) {
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(populateXRequestID),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := v1pb.RegisterJannaAPIHandlerFromEndpoint(ctx, gwMux, ln.Addr().String(), opts); err != nil {
		return nil, fmt.Errorf("failed to start HTTP gateway: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", contextWrap(gwMux))

	server := &http.Server{
		Handler: grpcHandlerFunc(grpcServer, mux),
	}

	return server, nil
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
