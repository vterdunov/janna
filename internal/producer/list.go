package producer

import "context"

// VMList is a command that implements a usecase that retrieve Virtual Machines list.
type VMList struct {
	params VMListRequest

	Producer
}

func NewVMList(params VMListRequest, producer Producer) *VMList {
	return &VMList{
		params:   params,
		Producer: producer,
	}
}

// Execute returns a Virtual Machine information
func (i *VMList) Execute(ctx context.Context) (VMListResponse, error) {
	return i.Producer.VMListTask(ctx, i.params)
}

type VMListRequest struct {
	Datacenter   string
	Folder       string
	ResourcePool string
}

type VMListResponse struct {
	TaskID string
}
