package usecase

func (u *Usecase) VMInfo(uuid string) (VMInfoResponse, error) {
	return u.vmWareRepository.VMInfo(uuid)
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
