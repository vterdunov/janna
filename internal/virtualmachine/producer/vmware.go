// producer produce a task for distributed workers. Each task puts into a queue and retruns the task ID.
package producer

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

type Producer struct {
	server *machinery.Server
}

func NewProducer(redisURL string) (*Producer, error) {
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

	w := Producer{
		server: server,
	}

	return &w, nil
}

func (w *Producer) VMDeployTask(ctx context.Context, params virtualmachine.VMDeployRequest) (virtualmachine.VMDeployResponse, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return virtualmachine.VMDeployResponse{}, err
	}

	sEnc := base64.StdEncoding.EncodeToString(buf.Bytes())
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
