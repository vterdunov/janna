package broker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/vterdunov/janna/internal/producer"
)

type RedisPublisher struct {
	server *machinery.Server
}

// NewRedisProducer creates implementation of producer.Producer interface
func NewRedisProducer(redisURL string) (RedisPublisher, error) {
	cnf := &config.Config{
		Broker:        redisURL,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: redisURL,
		Redis:         &config.RedisConfig{},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return RedisPublisher{}, err
	}

	redisPublisher := RedisPublisher{
		server: server,
	}

	return redisPublisher, nil
}

func (p RedisPublisher) TaskInfo(ctx context.Context, params producer.TaskInfoRequest) (producer.TaskInfoResponse, error) {
	asyncResult, err := p.server.GetBackend().GetState(params.TaskID)
	if err != nil {
		return producer.TaskInfoResponse{}, fmt.Errorf("could not get info for task id: %s", params.TaskID)
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(asyncResult.Results); err != nil {
		return producer.TaskInfoResponse{}, err
	}

	var data string
	for _, r := range asyncResult.Results {
		data = r.Value.(string)
	}

	var taskType producer.TaskType

	switch asyncResult.TaskName {
	case "vm_deploy":
		taskType = producer.VMDeployTask
	case "vm_info":
		taskType = producer.VMInfoTask
	default:
		taskType = producer.Invalid
	}

	result := producer.TaskInfoResponse{
		State:    asyncResult.State,
		TaskType: taskType,
		Data:     data,
		Err:      fmt.Errorf("%s", asyncResult.Error),
	}

	return result, nil
}

func (p RedisPublisher) VMDeployTask(ctx context.Context, params producer.VMDeployRequest) (producer.VMDeployResponse, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return producer.VMDeployResponse{}, err
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
		return producer.VMDeployResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	return producer.VMDeployResponse{
		TaskID: taskID,
	}, nil
}

func (p RedisPublisher) VMInfoTask(ctx context.Context, params producer.VMInfoRequest) (producer.VMInfoResponse, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return producer.VMInfoResponse{}, err
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
		return producer.VMInfoResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	return producer.VMInfoResponse{
		TaskID: taskID,
	}, nil
}

func (p RedisPublisher) VMListTask(ctx context.Context, params producer.VMListRequest) (producer.VMListResponse, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return producer.VMListResponse{}, err
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
		return producer.VMListResponse{}, err
	}

	taskID := asyncResult.GetState().TaskUUID

	return producer.VMListResponse{
		TaskID: taskID,
	}, nil
}
