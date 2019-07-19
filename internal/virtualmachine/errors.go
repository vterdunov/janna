package virtualmachine

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrVMAlreadyExist = Error("Virtual Machine already exist")
