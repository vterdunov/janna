package worker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/vterdunov/janna/internal/virtualmachine"
)

type Worker struct {
	server *machinery.Server
}

func NewWorker(redisURL string) (*Worker, error) {
	cnf := &config.Config{
		Broker:        redisURL,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: redisURL,
		Redis:         &config.RedisConfig{},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	w := Worker{
		server: server,
	}

	return &w, nil
}

func (w *Worker) VMDeploy(ctx context.Context, params virtualmachine.VMDeployRequest) (virtualmachine.VMDeployResponse, error) {
	// var crType string
	// switch params.ComputerResources.Type {
	// case virtualmachine.ComputerResourceHost:
	// 	crType = "host"
	// case virtualmachine.ComputerResourceCluster:
	// 	crType = "cluster"
	// case virtualmachine.ComputerResourceResourcePool:
	// 	crType = "rp"
	// }

	// cr := fmt.Sprintf("%v,%v", crType, params.ComputerResources.Path)

	// var dsType string
	// switch params.Datastores.Type {
	// case virtualmachine.DatastoreCluster:
	// 	dsType = "cluster"
	// case virtualmachine.DatastoreDatastore:
	// 	dsType = "datastore"
	// }

	// ds := fmt.Sprintf("%v,%v", dsType, params.Datastores.Names)

	// taskValues := []string{
	// 	"name", params.Name,
	// 	"datacenter", params.Datacenter,
	// 	"ova_url", params.OvaURL,
	// 	"folder", params.Folder,
	// 	"annotation", params.Annotation,
	// 	"computer_resources", cr,
	// 	"datastores", ds,
	// }

	// signature := &tasks.Signature{
	// 	Name: "vm_deploy",
	// 	Args: []tasks.Arg{
	// 		{
	// 			Type:  "[]string",
	// 			Value: taskValues,
	// 		},
	// 	},
	// }

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(params)
	if err != nil {
		return virtualmachine.VMDeployResponse{}, err
	}

	sEnc := base64.StdEncoding.EncodeToString(network.Bytes())
	signature := &tasks.Signature{
		Name: "vm_deploy",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: sEnc,
			},
		},
	}

	asyncResult, err := w.server.SendTask(signature)
	if err != nil {
		return virtualmachine.VMDeployResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	// taskState, err := server.GetBackend().GetState(taskState.TaskUUID)

	return virtualmachine.VMDeployResponse{
		TaskID: taskID,
	}, nil
}
