// usecase contains the application Use Cases
// it provides interaction between delivery, repository and domain layers
package usecase

type Usecase struct {
	appInfoRepository AppInfoRepository
	vmWareRepository  VMWareRepository
}

type AppInfoRepository interface {
	AppInfo() (*AppInfoResponse, error)
}

type VMWareRepository interface {
	VMInfo(uuid string) (VMInfoResponse, error)
	VMDeploy(params VMDeployRequest) (VMDeployResponse, error)
}

func NewUsecase(r AppInfoRepository, wmwareRep VMWareRepository) *Usecase {
	return &Usecase{
		appInfoRepository: r,
		vmWareRepository:  wmwareRep,
	}
}
