package usecase

type VMInfo struct {
	params VMInfoRequest
	VMRepository
}

func NewVMInfo() *VMInfo {
	return &VMInfo{}
}

func (i *VMInfo) Execute() (VMInfoResponse, error) {
	return i.vmInfo(i.params.uuid)
}

type VMInfoRequest struct {
	uuid string
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
