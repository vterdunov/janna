package virtualmachine

type VMRepository interface {
	VMInfo(uuid string) (VMInfoResponse, error)
	VMDeploy(params VMDeployRequest) (VMDeployResponse, error)
}
