package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/appinfo"
	"github.com/vterdunov/janna/internal/virtualmachine"
)

type server struct {
	appInfoRepository appinfo.Repository
	vmRepository      virtualmachine.VMRepository
}

func RegisterServer(
	gserver *grpc.Server,
	appInfoRepository appinfo.Repository,
	vmRepository virtualmachine.VMRepository) {
	s := &server{
		appInfoRepository: appInfoRepository,
		vmRepository:      vmRepository,
	}

	apiV1.RegisterJannaAPIServer(gserver, s)
	reflection.Register(gserver)
}

func (s *server) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	command := appinfo.NewAppInfo(s.appInfoRepository)

	appInfo := command.Execute()

	return &apiV1.AppInfoResponse{
		Commit:    appInfo.Commit,
		BuildTime: appInfo.BuildTime,
	}, nil
}

func (s *server) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	params := virtualmachine.VMInfoRequest{
		UUID: in.VmUuid,
	}

	command := virtualmachine.NewVMInfo(s.vmRepository, params)
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

	command := virtualmachine.NewVMDeploy(s.vmRepository, params)
	r, err := command.Execute()
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMDeployResponse{
		TaskId: r.TaskID,
	}

	return &resp, nil
}
