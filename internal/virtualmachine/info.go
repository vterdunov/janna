package virtualmachine

import "context"

// VMInfo is a command that implements a usecase that requests information about a Virtual Machine.
type VMInfo struct {
	params VMInfoRequest

	Producer
}

func NewVMInfo(params VMInfoRequest, producer Producer) *VMInfo {
	return &VMInfo{
		params:   params,
		Producer: producer,
	}
}

// Execute returns a Virtual Machine information
func (i *VMInfo) Execute(ctx context.Context) (VMInfoResponse, error) {
	return i.Producer.VMInfoTask(ctx, i.params)
}

type VMInfoRequest struct {
	UUID string
}

type VMInfoResponse struct {
	TaskID string
}
