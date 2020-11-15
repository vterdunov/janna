package middleware

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../../../../tools/templates/log template

//go:generate go run github.com/hexdigest/gowrap/cmd/gowrap gen -g -p github.com/vterdunov/janna-proto/gen/go/v1 -i JannaAPIServer -t ../../../../tools/templates/log -o logging_gen.go

import (
	context "context"
	"time"

	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	"github.com/pkg/errors"
	"github.com/vterdunov/janna/internal/log"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apiv1 "github.com/vterdunov/janna-proto/gen/go/v1"
)

// JannaAPIServerWithLog implements apiv1.JannaAPIServer that is instrumented with logrus logger
type JannaAPIServerWithLog struct {
	logger log.Logger
	next   apiv1.JannaAPIServer
}

// NewJannaAPIServerWithLog instruments an implementation of the apiv1.JannaAPIServer with simple logging
func NewJannaAPIServerWithLog(next apiv1.JannaAPIServer, logger log.Logger) JannaAPIServerWithLog {
	return JannaAPIServerWithLog{
		next:   next,
		logger: logger,
	}
}

// AppInfo implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) AppInfo(ctx context.Context, ap1 *apiv1.AppInfoRequest) (ap2 *apiv1.AppInfoResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "AppInfo",
	)

	logger.Info("calling endpoint")

	res, err := m.next.AppInfo(ctx, ap1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// OpenApi implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) OpenApi(ctx context.Context, op1 *apiv1.OpenApiRequest) (op2 *apiv1.OpenApiResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "OpenApi",
	)

	logger.Info("calling endpoint")

	res, err := m.next.OpenApi(ctx, op1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// TaskStatus implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) TaskStatus(ctx context.Context, tp1 *apiv1.TaskStatusRequest) (tp2 *apiv1.TaskStatusResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "TaskStatus",
	)

	logger.Info("calling endpoint")

	res, err := m.next.TaskStatus(ctx, tp1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// VMDeploy implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) VMDeploy(ctx context.Context, vp1 *apiv1.VMDeployRequest) (vp2 *apiv1.VMDeployResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMDeploy",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMDeploy(ctx, vp1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// VMInfo implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) VMInfo(ctx context.Context, vp1 *apiv1.VMInfoRequest) (vp2 *apiv1.VMInfoResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMInfo",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMInfo(ctx, vp1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// VMList implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) VMList(ctx context.Context, vp1 *apiv1.VMListRequest) (vp2 *apiv1.VMListResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMList",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMList(ctx, vp1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

// VMPower implements apiv1.JannaAPIServer
func (m JannaAPIServerWithLog) VMPower(ctx context.Context, vp1 *apiv1.VMPowerRequest) (vp2 *apiv1.VMPowerResponse, err error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMPower",
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMPower(ctx, vp1)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)
		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call successful")
		}
	}(begin)

	return res, translateError(err)
}

func withRequestID(ctx context.Context, logger log.Logger) log.Logger {
	reqID := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["request_id"]) > 0 {
			reqID = md["request_id"][0]
		}
	}

	l := logger.WithFields(
		"request_id", reqID,
	)

	return l
}

type statusError interface {
	GRPCStatus() *status.Status
}

func isGrpcStatusError(err error) bool {
	_, ok := err.(statusError)
	return ok
}

// translateError translate business logic erros to transport level errors.
// May become clumsy because of the need to check all exposed errors
// from all packages with business logic.
func translateError(err error) error {
	if err == nil {
		return nil
	}

	if isGrpcStatusError(err) {
		return err
	}

	switch errors.Cause(err) {
	// case producer.ErrVMAlreadyExist:
	// 	err = status.Errorf(codes.AlreadyExists, err.Error())
	default:
		err = status.Errorf(codes.Internal, err.Error())
	}

	return err
}
