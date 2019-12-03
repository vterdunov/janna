package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/log"
	"github.com/vterdunov/janna/internal/virtualmachine"
)

// Service implements apiV1.JannaAPIServer
type Service struct {
	appInfoRepository appinfo.Repository
	producer          virtualmachine.Producer
}

func NewService(i appinfo.Repository, p virtualmachine.Producer) apiV1.JannaAPIServer {
	return Service{
		appInfoRepository: i,
		producer:          p,
	}
}

func RegisterServer(gserver *grpc.Server, service apiV1.JannaAPIServer, logger log.Logger) {
	apiV1.RegisterJannaAPIServer(gserver, service)
	reflection.Register(gserver)
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
	params := virtualmachine.VMInfoRequest{
		UUID: in.VmUuid,
	}

	command := virtualmachine.NewVMInfo(params, s.producer)
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
	var crType virtualmachine.ComputerResourcesType
	var crPath string
	if in.ComputerResources != nil {
		crPath = in.ComputerResources.Path

		switch in.ComputerResources.Type.String() {
		case "TYPE_HOST":
			crType = virtualmachine.ComputerResourceHost
		case "TYPE_CLUSTER":
			crType = virtualmachine.ComputerResourceCluster
		case "TYPE_RP":
			crType = virtualmachine.ComputerResourceResourcePool
		default:
			return nil, errors.New("could not recognize Computer resource type. Please read documentation")
		}
	}

	cr := virtualmachine.ComputerResources{
		Type: crType,
		Path: crPath,
	}

	var dsType virtualmachine.DatastoreType
	var dsNames []string
	if in.Datastores != nil {
		dsNames = in.Datastores.Names

		switch in.Datastores.Type.String() {
		case "TYPE_CLUSTER":
			dsType = virtualmachine.DatastoreCluster
		case "TYPE_DATASTORE":
			dsType = virtualmachine.DatastoreDatastore
		default:
			return nil, errors.New("could not recognize Datastore type. Please read documentation")
		}
	}

	datastores := virtualmachine.Datastores{
		Type:  dsType,
		Names: dsNames,
	}

	params := virtualmachine.VMDeployRequest{
		Name:              in.Name,
		Datacenter:        in.Datacenter,
		OvaURL:            in.OvaUrl,
		Folder:            in.Folder,
		Annotation:        in.Annotation,
		ComputerResources: cr,
		Datastores:        datastores,
	}

	command := virtualmachine.NewVMDeploy(params, s.producer)
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
	params := virtualmachine.VMListRequest{
		Datacenter:   in.Datacenter,
		Folder:       in.Folder,
		ResourcePool: in.ResourcePool,
	}

	command := virtualmachine.NewVMList(params, s.producer)
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
