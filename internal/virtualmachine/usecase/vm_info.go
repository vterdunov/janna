package usecase

type VMInfo struct {
	VMRepository
}

func (i *VMInfo) New() VMInfo {
	return VMInfo{}
}

func (i *VMInfo) Execute(uuid string) (VMInfoResponse, error) {
	return i.VMInfo(uuid)
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
