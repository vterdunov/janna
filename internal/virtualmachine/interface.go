package virtualmachine

// VMRepository describes abstract methods to interact with VMWare servers in Virtual Machine bounded context.
type VMRepository interface {
	VMInfo(uuid string) (VMInfoResponse, error)
	VMDeploy(params VMDeployRequest) (VMDeployResponse, error)
}
