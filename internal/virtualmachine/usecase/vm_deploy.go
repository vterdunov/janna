package usecase

type DatastoreType int

const (
	// Datastore Type Enum
	DatastoreInvalid DatastoreType = iota
	DatastoreCluster
	DatastoreDatastore
)

type ComputerResourcesType int

const (
	// Computer Resources Enum
	ComputerResourceInvalid ComputerResourcesType = iota
	ComputerResourceHost
	ComputerResourceCluster
	ComputerResourceResourcePool
)

type VMDeploy struct {
	VMRepository
}

func (d *VMDeploy) New() VMDeploy {
	return VMDeploy{}
}

func (d *VMDeploy) Execute(params VMDeployRequest) (VMDeployResponse, error) {
	return d.VMDeploy(params)
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
	Type ComputerResourcesType
}

type Datastores struct {
	Type  DatastoreType
	Names []string
}
