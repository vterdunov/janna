package grpc

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	stpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/log"
	"github.com/vterdunov/janna/internal/producer"
)

// Service implements apiV1.JannaAPIServer
type Service struct {
	appInfoRepository appinfo.Repository
	producer          producer.Producer
}

func NewService(i appinfo.Repository, p producer.Producer) apiV1.JannaAPIServer {
	return Service{
		appInfoRepository: i,
		producer:          p,
	}
}

func RegisterServer(gserver *grpc.Server, service apiV1.JannaAPIServer, logger log.Logger) {
	apiV1.RegisterJannaAPIServer(gserver, service)
	reflection.Register(gserver)
}

func (s Service) TaskStatus(ctx context.Context, in *apiV1.TaskStatusRequest) (*apiV1.TaskStatusResponse, error) {
	params := producer.TaskInfoRequest{
		TaskID: in.TaskId,
	}

	command := producer.NewTaskInfo(params, s.producer)
	result, err := command.Execute(ctx)
	if err != nil {
		return nil, err
	}

	dResult, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return nil, fmt.Errorf("could not decode result: %w", err)
	}

	item := stpb.Struct{
		Fields: map[string]*stpb.Value{
			"testkey": &stpb.Value{
				Kind: &stpb.Value_StringValue{StringValue: "test string"},
			},
		},
	}

	resp := apiV1.TaskStatusResponse{
		Status:  "test status",
		Message: "test message",
		Result:  &item,
	}

	return &resp, nil
}

func (s Service) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	command := appinfo.NewAppInfo(s.appInfoRepository)

	appInfo := command.Execute()

	return &apiV1.AppInfoResponse{
		Commit:    appInfo.Commit,
		BuildTime: appInfo.BuildTime,
	}, nil
}

func (s Service) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	params := producer.VMInfoRequest{
		UUID: in.VmUuid,
	}

	command := producer.NewVMInfo(params, s.producer)
	info, err := command.Execute(ctx)
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMInfoResponse{
		TaskId: info.TaskID,
	}

	return &resp, nil
}

func (s Service) VMDeploy(ctx context.Context, in *apiV1.VMDeployRequest) (*apiV1.VMDeployResponse, error) {
	// TODO: validate incoming data
	var crType producer.ComputerResourcesType
	var crPath string
	if in.ComputerResources != nil {
		crPath = in.ComputerResources.Path

		switch in.ComputerResources.Type.String() {
		case "TYPE_HOST":
			crType = producer.ComputerResourceHost
		case "TYPE_CLUSTER":
			crType = producer.ComputerResourceCluster
		case "TYPE_RP":
			crType = producer.ComputerResourceResourcePool
		default:
			return nil, errors.New("could not recognize Computer resource type. Please read documentation")
		}
	}

	cr := producer.ComputerResources{
		Type: crType,
		Path: crPath,
	}

	var dsType producer.DatastoreType
	var dsNames []string
	if in.Datastores != nil {
		dsNames = in.Datastores.Names

		switch in.Datastores.Type.String() {
		case "TYPE_CLUSTER":
			dsType = producer.DatastoreCluster
		case "TYPE_DATASTORE":
			dsType = producer.DatastoreDatastore
		default:
			return nil, errors.New("could not recognize Datastore type. Please read documentation")
		}
	}

	datastores := producer.Datastores{
		Type:  dsType,
		Names: dsNames,
	}

	params := producer.VMDeployRequest{
		Name:              in.Name,
		Datacenter:        in.Datacenter,
		OvaURL:            in.OvaUrl,
		Folder:            in.Folder,
		Annotation:        in.Annotation,
		ComputerResources: cr,
		Datastores:        datastores,
	}

	command := producer.NewVMDeploy(params, s.producer)
	r, err := command.Execute(ctx)
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMDeployResponse{
		TaskId: r.TaskID,
	}

	return &resp, nil
}

func (s Service) VMList(ctx context.Context, in *apiV1.VMListRequest) (*apiV1.VMListResponse, error) {
	params := producer.VMListRequest{
		Datacenter:   in.Datacenter,
		Folder:       in.Folder,
		ResourcePool: in.ResourcePool,
	}

	command := producer.NewVMList(params, s.producer)
	r, err := command.Execute(ctx)
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMListResponse{
		TaskId: r.TaskID,
	}
	return &resp, nil
}

func (s Service) VMPower(ctx context.Context, in *apiV1.VMPowerRequest) (*apiV1.VMPowerResponse, error) {
	return nil, errors.New("not implemented")
}
