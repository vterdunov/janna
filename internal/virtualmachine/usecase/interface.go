package usecase

type VMRepository interface {
	vmInfo(uuid string) (VMInfoResponse, error)
	vmDeploy(params VMDeployRequest) (VMDeployResponse, error)
}
