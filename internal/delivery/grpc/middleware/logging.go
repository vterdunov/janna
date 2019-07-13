package middleware

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/log"
)

func NewMiddleware(next apiV1.JannaAPIServer, logger log.Logger) apiV1.JannaAPIServer {
	service := ErrorHandlingMiddleware{
		logger: logger,
		next:   next,
	}

	return &service
}

type ErrorHandlingMiddleware struct {
	logger log.Logger
	next   apiV1.JannaAPIServer
}

func (m *ErrorHandlingMiddleware) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "AppInfo",
	)

	logger.Info("calling endpoint")

	res, err := m.next.AppInfo(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "called endpoint")
		} else {
			logger.Info("called endpoint")
		}

	}(begin)

	return res, err
}

func (m *ErrorHandlingMiddleware) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMInfo",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMInfo(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "called endpoint")
		} else {
			logger.Info("called endpoint")
		}

	}(begin)

	return res, err
}

func (m *ErrorHandlingMiddleware) VMDeploy(ctx context.Context, in *apiV1.VMDeployRequest) (*apiV1.VMDeployResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMDeploy",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMDeploy(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "called endpoint")
		} else {
			logger.Info("called endpoint")
		}

	}(begin)

	return res, err
}

func (m *ErrorHandlingMiddleware) VMList(ctx context.Context, in *apiV1.VMListRequest) (*apiV1.VMListResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMList",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMList(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "called endpoint")
		} else {
			logger.Info("called endpoint")
		}

	}(begin)

	return res, err
}

func (m *ErrorHandlingMiddleware) VMPower(ctx context.Context, in *apiV1.VMPowerRequest) (*apiV1.VMPowerResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMPower",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMPower(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "called endpoint")
		} else {
			logger.Info("called endpoint")
		}

	}(begin)

	return res, err
}

func withRequestID(ctx context.Context, logger log.Logger) log.Logger {
	reqID := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqID = md["request_id"][0]
	}

	l := logger.WithFields(
		"request_id", reqID,
	)

	return l
}
