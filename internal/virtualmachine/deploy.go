package virtualmachine

import (
	"context"
)

type DatastoreType int

const (
	// Datastore Type Enum
	DatastoreInvalid DatastoreType = iota
	DatastoreCluster
	DatastoreDatastore
)

type ComputerResourcesType int

const (
	// Computer Resources Enum
	ComputerResourceInvalid ComputerResourcesType = iota
	ComputerResourceHost
	ComputerResourceCluster
	ComputerResourceResourcePool
)

// VMDeploy is a command that implements a usecase that deploy a Virtual Machine from OVA file.
type VMDeploy struct {
	params VMDeployRequest

	Worker
	VMRepository
}

func NewVMDeploy(r VMRepository, params VMDeployRequest, worker Worker) VMDeploy {
	return VMDeploy{
		params:       params,
		VMRepository: r,
		Worker:       worker,
	}
}

func (d *VMDeploy) Execute(ctx context.Context) (VMDeployResponse, error) {
	// TODO: for speed up testing it was commented. Uncommend after distrubuted task was implemented
	// exist, err := d.IsVMExist(ctx, d.params.Name, d.params.Datacenter)
	// if err != nil {
	// 	return VMDeployResponse{}, err
	// }

	// if exist {
	// 	return VMDeployResponse{}, ErrVMAlreadyExist
	// }

	return d.VMDeploy(ctx, d.params)
}

type VMDeployRequest struct {
	Name       string
	Datacenter string
	OvaURL     string
	Folder     string
	Annotation string

	ComputerResources
	Datastores
}

type VMDeployResponse struct {
	TaskID string
}

type ComputerResources struct {
	Path string
	Type ComputerResourcesType
}

type Datastores struct {
	Type  DatastoreType
	Names []string
}
