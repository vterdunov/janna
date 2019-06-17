package usecase

type AppInfoRepository interface {
	AppInfo() (AppInfoResponse, error)
}

type VMRepository interface {
	VMInfo(uuid string) (VMInfoResponse, error)
	VMDeploy(params VMDeployRequest) (VMDeployResponse, error)
}
