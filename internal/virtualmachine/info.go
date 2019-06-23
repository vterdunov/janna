package virtualmachine

// VMInfo is a command that implements a usecase that requests information about a Virtual Machine.
type VMInfo struct {
	params VMInfoRequest
	VMRepository
}

func NewVMInfo(r VMRepository, params VMInfoRequest) *VMInfo {
	return &VMInfo{
		params:       params,
		VMRepository: r,
	}
}

// Execute returns a Virtual Machine information
func (i *VMInfo) Execute() (VMInfoResponse, error) {
	return i.VMInfo(i.params.UUID)
}

type VMInfoRequest struct {
	UUID string
}

type VMInfoResponse struct {
	Name             string
	UUID             string
	GuestID          string
	Annotation       string
	PowerState       string
	NumCPU           uint32
	NumEthernetCards uint32
	NumVirtualDisks  uint32
	Template         bool
	IPs              []string
}
