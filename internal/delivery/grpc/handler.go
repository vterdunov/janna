package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/virtualmachine/usecase"
)

type server struct {
	appInfo *usecase.AppInfo
	vmInfo *usecase.VMInfo
	vmDeploy *usecase.VMDeploy
}

func RegisterServer(
	gserver *grpc.Server,
	appInfo *usecase.AppInfo,
	vmInfo *usecase.VMInfo,
	vmDeploy *usecase.VMDeploy) {
	s := &server{
		appInfo: appInfo,
		vmInfo: vmInfo,
		vmDeploy: vmDeploy,
	}

	apiV1.RegisterJannaAPIServer(gserver, s)
	reflection.Register(gserver)
}

func (s *server) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	appInfo, err := s.appInfo.AppInfo()
	if err != nil {
		return nil, err
	}

	return &apiV1.AppInfoResponse{
		Commit:    appInfo.Commit,
		BuildTime: appInfo.BuildTime,
	}, nil
}

func (s *server) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	info, err := s.vmInfo.VMInfo(in.VmUuid)
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
	var crType usecase.ComputerResourcesType
	var crPath string
	if in.ComputerResources != nil {
		crPath = in.ComputerResources.Path

		switch in.ComputerResources.Type.String() {
		case "TYPE_HOST":
			crType = usecase.ComputerResourceHost
		case "TYPE_CLUSTER":
			crType = usecase.ComputerResourceCluster
		case "TYPE_RP":
			crType = usecase.ComputerResourceResourcePool
		default:
			return nil, errors.New("could not recognize Computer resource type. Please read documentation")
		}
	}

	cr := usecase.ComputerResources{
		Type: crType,
		Path: crPath,
	}

	var dsType usecase.DatastoreType
	var dsNames []string
	if in.Datastores != nil {
		dsNames = in.Datastores.Names

		switch in.Datastores.Type.String() {
		case "TYPE_CLUSTER":
			dsType = usecase.DatastoreCluster
		case "TYPE_DATASTORE":
			dsType = usecase.DatastoreDatastore
		default:
			return nil, errors.New("could not recognize Datastore type. Please read documentation")
		}
	}

	datastores := usecase.Datastores{
		Type:  dsType,
		Names: dsNames,
	}

	params := usecase.VMDeployRequest{
		Name:              in.Name,
		Datacenter:        in.Datacenter,
		OvaURL:            in.OvaUrl,
		Folder:            in.Folder,
		Annotation:        in.Annotation,
		ComputerResources: cr,
		Datastores:        datastores,
	}

	r, err := s.vmDeploy.VMDeploy(params)
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMDeployResponse{
		TaskId: r.TaskID,
	}

	return &resp, nil
}
