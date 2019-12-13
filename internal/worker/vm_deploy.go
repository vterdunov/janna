package worker

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

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

type VMDeploy struct{}

func (d *VMDeploy) Execute(params string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return "", fmt.Errorf("could not decode parameters from base64: %w", err)
	}

	r := bytes.NewReader(sDec)
	dec := gob.NewDecoder(r)

	var deployParams VMDeployRequest
	err = dec.Decode(&deployParams)
	if err != nil {
		return "", fmt.Errorf("could not decode parameters from bytes: %w", err)
	}

	spew.Dump(deployParams)

	dummyResp := "dummy dummy dummy dummy dummy dummy dummy dummy"

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(dummyResp); err != nil {
		return "", err
	}

	resp := base64.StdEncoding.EncodeToString(buf.Bytes())

	return resp, nil
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

type ComputerResources struct {
	Path string
	Type ComputerResourcesType
}

type Datastores struct {
	Type  DatastoreType
	Names []string
}

type VMDeployResponse struct {
	Name string
	IP   string
}
