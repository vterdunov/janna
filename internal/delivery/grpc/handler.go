package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/vterdunov/janna-proto/gen/go/v1"
	"github.com/vterdunov/janna/internal/usecase"
)

type server struct {
	usecase *usecase.Usecase
}

func RegisterServer(gserver *grpc.Server, u *usecase.Usecase) {
	s := &server{
		usecase: u,
	}

	apiV1.RegisterJannaAPIServer(gserver, s)
	reflection.Register(gserver)
}

func (s *server) AppInfo(ctx context.Context, in *apiV1.AppInfoRequest) (*apiV1.AppInfoResponse, error) {
	appInfo, err := s.usecase.AppInfo()
	if err != nil {
		return nil, err
	}

	return &apiV1.AppInfoResponse{
		Commit:    appInfo.Commit,
		BuildTime: appInfo.BuildTime,
	}, nil
}

func (s *server) VMInfo(ctx context.Context, in *apiV1.VMInfoRequest) (*apiV1.VMInfoResponse, error) {
	info, err := s.usecase.VMInfo(in.VmUuid)
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
	var crType string
	var crPath string
	if in.ComputerResources != nil {
		crType = in.ComputerResources.Type.String()
		crPath = in.ComputerResources.Path
	}

	cr := usecase.ComputerResources{
		Type: crType,
		Path: crPath,
	}

	var dsType string
	var dsNames []string
	if in.Datastores != nil {
		dsType = in.Datastores.Type.String()
		dsNames = in.Datastores.Names
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

	r, err := s.usecase.VMDeploy(params)
	if err != nil {
		return nil, err
	}

	resp := apiV1.VMDeployResponse{
		TaskId: r.TaskID,
	}

	return &resp, nil
}
