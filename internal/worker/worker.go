package worker

import (
	"errors"
)

type VMDeploy struct {
}

func (d *VMDeploy) Execute(params string) error {
	return errors.New("VMDeploy not inplemented")
}

type VMList struct {
}

func (d *VMList) Execute(params string) error {
	return errors.New("VMList not inplemented")
}

type VMInfo struct {
}

func (d *VMInfo) Execute(params string) error {
	return errors.New("VMInfo not inplemented")
}
