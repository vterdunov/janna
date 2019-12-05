package middleware

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/log"
	"github.com/vterdunov/janna/internal/producer"
)

func NewLoggingMiddleware(next apiV1.JannaAPIServer, logger log.Logger) apiV1.JannaAPIServer {
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
			logger.Error(err, "call failed")
		} else {
			logger.Info("call finished")
		}

	}(begin)

	return res, translateError(err)
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
			logger.Error(err, "call failed")
		} else {
			logger.Info("call finished")
		}

	}(begin)

	return res, translateError(err)
}

func (m *ErrorHandlingMiddleware) VMDeploy(ctx context.Context, in *apiV1.VMDeployRequest) (*apiV1.VMDeployResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMDeploy",
		"vm_name", in.Name,
		"ova_url", in.OvaUrl,
		"datacenter", in.Datacenter,
		"folder", in.Folder,
		"annotation", in.Annotation,
		"networks", in.Networks,
		"datastores", in.Datastores.String(),
		"computer_resources", in.ComputerResources.String(),
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMDeploy(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call finished")
		}

	}(begin)

	return res, translateError(err)
}

func (m *ErrorHandlingMiddleware) VMList(ctx context.Context, in *apiV1.VMListRequest) (*apiV1.VMListResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMList",
		"datacenter", in.Datacenter,
		"folder", in.Folder,
		"resource_pool", in.ResourcePool,
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMList(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call finished")
		}

	}(begin)

	return res, translateError(err)
}

func (m *ErrorHandlingMiddleware) VMPower(ctx context.Context, in *apiV1.VMPowerRequest) (*apiV1.VMPowerResponse, error) {
	begin := time.Now()
	logger := withRequestID(ctx, m.logger)
	logger = logger.WithFields(
		"method", "VMPower",
		"vm_uuid", in.VmUuid,
		"vm_power_request_body", in.VmPowerRequestBody.String(),
	)

	logger.Info("calling endpoint")

	res, err := m.next.VMPower(ctx, in)

	defer func(begin time.Time) {
		logger = logger.WithFields(
			"took", time.Since(begin).String(),
		)

		if err != nil {
			logger.Error(err, "call failed")
		} else {
			logger.Info("call finished")
		}

	}(begin)

	return res, translateError(err)
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

type statusError interface {
	GRPCStatus() *status.Status
}

func isGrpcStatusError(err error) bool {
	_, ok := err.(statusError)
	return ok
}

// translateError translate business logic erros to transport level errors.
// May become a clumsy because the need to check all exposed errors
// from all packages with business logic.
func translateError(err error) error {
	if err == nil {
		return nil
	}

	if isGrpcStatusError(err) {
		return err
	}

	switch errors.Cause(err) {
	case producer.ErrVMAlreadyExist:
		err = status.Errorf(codes.AlreadyExists, err.Error())
	default:
		err = status.Errorf(codes.Internal, err.Error())
	}

	return err
}
