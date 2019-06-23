package virtualmachine

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

// VMDeploy is a command that implements a usecase that deploy a Virtual Machine from OVA file.
type VMDeploy struct {
	params VMDeployRequest
	VMRepository
}

func NewVMDeploy(r VMRepository, params VMDeployRequest) VMDeploy {
	return VMDeploy{
		params:       params,
		VMRepository: r,
	}
}

func (d *VMDeploy) Execute() (VMDeployResponse, error) {
	return d.VMDeploy(d.params)
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