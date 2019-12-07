// producer produce a task for distributed workers. Each task puts into a queue and retruns the task ID.
package producer

import "context"

// Producer describes some producer that can sends tasks to distributed workers
type Producer interface {
	VMDeployTask(context.Context, VMDeployRequest) (VMDeployResponse, error)
	VMInfoTask(context.Context, VMInfoRequest) (VMInfoResponse, error)
	VMListTask(context.Context, VMListRequest) (VMListResponse, error)
}
