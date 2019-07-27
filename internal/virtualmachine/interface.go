package virtualmachine

import "context"

// VMRepository describes abstract methods to interact with VMWare servers in Virtual Machine bounded context.
type VMRepository interface {
	IsVMExist(context.Context, string, string) (bool, error)
	VMList(VMListRequest) ([]VMListResponse, error)
	VMInfo(uuid string) (VMInfoResponse, error)
}

// Worker describes some worker that can perform (often) long-running tasks
type Worker interface {
	VMDeploy(context.Context, VMDeployRequest) (VMDeployResponse, error)
}
