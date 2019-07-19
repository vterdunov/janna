package virtualmachine

import "context"

// VMRepository describes abstract methods to interact with VMWare servers in Virtual Machine bounded context.
type VMRepository interface {
	IsVMExist(context.Context, string, string) (bool, error)
	VMList(VMListRequest) ([]VMListResponse, error)
	VMInfo(uuid string) (VMInfoResponse, error)
	VMDeploy(context.Context, VMDeployRequest) (VMDeployResponse, error)
}

// StatusStorager represents behavior of storage that keeps deploy jobs statuses
type StatusStorager interface {
	NewTask() TaskStatuser
	FindByID(id string) TaskStatuser
}

// TaskStatuser represents behavior of every single task
type TaskStatuser interface {
	ID() string
	Str(keyvals ...string) TaskStatuser
	StrArr(key string, arr []string) TaskStatuser
	Get() (statuses map[string]interface{})
}
