package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/appinfo"
	vmRepository "github.com/vterdunov/janna/internal/virtualmachine/repository"
	vmUsecase "github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

type server struct {
	appInfoRepository appinfo.AppInfoRepository
	vmRepository      vmRepository.VMRepository
}

func RegisterServer(
	gserver *grpc.Server,
	appInfoRepository appinfo.AppInfoRepository,
	vmRepository vmRepository.VMRepository) {
	s := &server{
		appInfoRepository: appInfoRepository,
		vmRepository:      vmRepository,
	}

	apiV1.RegisterJannaAPIServer(gserver, s)
	reflection.Register(gserver)
}

func (s *server) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	command := appinfo.NewAppInfo(s.appInfoRepository)

	appInfo, err := command.Execute()
	if err != nil {
		return nil, err
	}

	return &apiV1.AppInfoResponse{
		Commit:    appInfo.Commit,
		BuildTime: appInfo.BuildTime,
	}, nil
}

func (s *server) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	command := vmUsecase.NewVMInfo()
	info, err := command.Execute()
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMInfoResponse{
		Name:             info.Name,
		Uuid:             info.UUID,
		GuestId:          info.GuestID,
		Annotation:       info.Annotation,
		PowerState:       info.PowerState,
		NumCpu:           info.NumCPU,
		NumEthernetCards: info.NumEthernetCards,
		NumVirtualDisks:  info.NumVirtualDisks,
		Template:         info.Template,
	}

	return &resp, nil
}

func (s *server) VMDeploy(ctx context.Context, in *apiV1.VMDeployRequest) (*apiV1.VMDeployResponse, error) {
	// TODO: validate incoming data
	var crType vmUsecase.ComputerResourcesType
	var crPath string
	if in.ComputerResources != nil {
		crPath = in.ComputerResources.Path

		switch in.ComputerResources.Type.String() {
		case "TYPE_HOST":
			crType = vmUsecase.ComputerResourceHost
		case "TYPE_CLUSTER":
			crType = vmUsecase.ComputerResourceCluster
		case "TYPE_RP":
			crType = vmUsecase.ComputerResourceResourcePool
		default:
			return nil, errors.New("could not recognize Computer resource type. Please read documentation")
		}
	}

	cr := vmUsecase.ComputerResources{
		Type: crType,
		Path: crPath,
	}

	var dsType vmUsecase.DatastoreType
	var dsNames []string
	if in.Datastores != nil {
		dsNames = in.Datastores.Names

		switch in.Datastores.Type.String() {
		case "TYPE_CLUSTER":
			dsType = vmUsecase.DatastoreCluster
		case "TYPE_DATASTORE":
			dsType = vmUsecase.DatastoreDatastore
		default:
			return nil, errors.New("could not recognize Datastore type. Please read documentation")
		}
	}

	datastores := vmUsecase.Datastores{
		Type:  dsType,
		Names: dsNames,
	}

	params := vmUsecase.VMDeployRequest{
		Name:              in.Name,
		Datacenter:        in.Datacenter,
		OvaURL:            in.OvaUrl,
		Folder:            in.Folder,
		Annotation:        in.Annotation,
		ComputerResources: cr,
		Datastores:        datastores,
	}

	command := vmUsecase.NewVMDeploy(s.vmRepository, params)
	r, err := command.Execute()
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMDeployResponse{
		TaskId: r.TaskID,
	}

	return &resp, nil
}
