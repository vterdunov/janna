package producer

import "context"

type TaskInfo struct {
	params TaskInfoRequest

	Producer
}

func NewTaskInfo(params TaskInfoRequest, producer Producer) TaskInfo {
	return TaskInfo{
		params:   params,
		Producer: producer,
	}
}

func (d *TaskInfo) Execute(ctx context.Context) (string, error) {
	return d.TaskInfo(ctx, d.params)
}

type TaskInfoRequest struct {
	TaskID string
}
