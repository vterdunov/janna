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
)

type Publisher struct {
	server *machinery.Server
}

func NewProducer(redisURL string) (*Publisher, error) {
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

	w := Publisher{
		server: server,
	}

	return &w, nil
}

func (p *Publisher) VMDeployTask(ctx context.Context, params VMDeployRequest) (VMDeployResponse, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return VMDeployResponse{}, err
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

	asyncResult, err := p.server.SendTask(signature)
	if err != nil {
		return VMDeployResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	// taskState, err := server.GetBackend().GetState(taskState.TaskUUID)

	return VMDeployResponse{
		TaskID: taskID,
	}, nil
}

func (p *Publisher) VMInfoTask(ctx context.Context, params VMInfoRequest) (VMInfoResponse, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return VMInfoResponse{}, err
	}

	sEnc := base64.StdEncoding.EncodeToString(buf.Bytes())
	signature := &tasks.Signature{
		Name: "vm_info",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: sEnc,
			},
		},
	}

	asyncResult, err := p.server.SendTask(signature)
	if err != nil {
		return VMInfoResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	// taskState, err := server.GetBackend().GetState(taskState.TaskUUID)

	return VMInfoResponse{
		TaskID: taskID,
	}, nil
}

func (p *Publisher) VMListTask(ctx context.Context, params VMListRequest) (VMListResponse, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return VMListResponse{}, err
	}

	sEnc := base64.StdEncoding.EncodeToString(buf.Bytes())
	signature := &tasks.Signature{
		Name: "vm_list",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: sEnc,
			},
		},
	}

	asyncResult, err := p.server.SendTask(signature)
	if err != nil {
		return VMListResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	// taskState, err := server.GetBackend().GetState(taskState.TaskUUID)

	return VMListResponse{
		TaskID: taskID,
	}, nil
}
