package virtualmachine

// VMList is a command that implements a usecase that retreive Virtual Machines list.
type VMList struct {
	params VMListRequest
	VMRepository
}

func NewVMList(r VMRepository, params VMListRequest) *VMList {
	return &VMList{
		params:       params,
		VMRepository: r,
	}
}

// Execute returns a Virtual Machine information
func (i *VMList) Execute() ([]VMListResponse, error) {
	return i.VMList(i.params)
}

type VMListRequest struct {
	Datacenter   string
	Folder       string
	ResourcePool string
}

type VMListResponse struct {
	Name string
	UUID string
}
