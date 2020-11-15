package worker

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type VMList struct {
}

func (d *VMList) Execute(params string) error {
	sDec, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return fmt.Errorf("could not decode parameters from base64: %w", err)
	}

	r := bytes.NewReader(sDec)
	dec := gob.NewDecoder(r)

	var listParams VMListRequest
	err = dec.Decode(&listParams)
	if err != nil {
		return fmt.Errorf("could not decode parameters from bytes: %w", err)
	}

	spew.Dump(listParams)
	return nil
}

type VMListRequest struct {
	Datacenter   string
	Folder       string
	ResourcePool string
}
