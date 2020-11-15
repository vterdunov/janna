package worker

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type VMInfo struct {
}

func (d *VMInfo) Execute(params string) error {
	sDec, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return fmt.Errorf("could not decode parameters from base64: %w", err)
	}

	r := bytes.NewReader(sDec)
	dec := gob.NewDecoder(r)

	var infoParams VMInfoRequest
	err = dec.Decode(&infoParams)
	if err != nil {
		return fmt.Errorf("could not decode parameters from bytes: %w", err)
	}

	spew.Dump(infoParams)
	return nil
}

type VMInfoRequest struct {
	UUID string
}
