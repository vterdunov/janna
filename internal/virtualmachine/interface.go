package virtualmachine

import "context"

// VMRepository describes abstract methods to interact with VMWare servers in Virtual Machine bounded context.
type VMRepository interface {
	IsVMExist(context.Context, string, string) (bool, error)
	VMList(VMListRequest) ([]VMListResponse, error)
	VMInfo(uuid string) (VMInfoResponse, error)
}

// Producer describes some producer that can sends tasks to distributed workers
type Producer interface {
	VMDeployTask(context.Context, VMDeployRequest) (VMDeployResponse, error)
	VMInfoTask(context.Context, VMInfoRequest) (VMInfoResponse, error)
	VMListTask(context.Context, VMListRequest) (VMListResponse, error)
}
