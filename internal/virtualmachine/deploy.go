package virtualmachine

import (
	"context"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
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

	StatusStorager
	VMRepository
}

func NewVMDeploy(r VMRepository, params VMDeployRequest, statusStorage StatusStorager) VMDeploy {
	return VMDeploy{
		params:         params,
		VMRepository:   r,
		StatusStorager: statusStorage,
	}
}

func (d *VMDeploy) Execute(ctx context.Context) (VMDeployResponse, error) {
	// exist, err := d.IsVMExist(ctx, d.params.Name, d.params.Datacenter)
	// if err != nil {
	// 	return VMDeployResponse{}, err
	// }

	// if exist {
	// 	return VMDeployResponse{}, ErrVMAlreadyExist
	// }

	t := d.NewTask()
	t.Str("stage", "start")

	cnf := &config.Config{
		Broker:        "redis://redis:6379",
		DefaultQueue:  "machinery_tasks",
		ResultBackend: "redis://redis:6379",
		Redis:         &config.RedisConfig{},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		log.Printf("err %v", err.Error())
		return VMDeployResponse{}, err
	}

	// server.RegisterTask("add", Add)

	signature := &tasks.Signature{
		Name: "add",
		Args: []tasks.Arg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 1,
			},
		},
	}

	asyncResult, err := server.SendTask(signature)
	if err != nil {
		log.Printf("err %v", err.Error())
		return VMDeployResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID
	_ = taskID

	// taskState, err := server.GetBackend().GetState(taskState.TaskUUID)

	// temporary. for testing purposes
	return VMDeployResponse{}, ErrVMAlreadyExist
	// return d.VMDeploy(ctx, d.params)
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
