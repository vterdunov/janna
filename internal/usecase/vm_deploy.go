package usecase

func (u *Usecase) VMDeploy(params VMDeployRequest) (VMDeployResponse, error) {
	return u.vmWareRepository.VMDeploy(params)
}

type VMDeployRequest struct {
	Name       string
	Datacenter string
	OvaURL     string
	Folder     string
	Annotation string

	ComputerResources
	Datastores
}

type VMDeployResponse struct {
	TaskID string
}

type ComputerResources struct {
	Path string
	Type string
}

type Datastores struct {
	Type  string
	Names []string
}
