package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func main() {
	cnf := &config.Config{
		Broker:        "redis://redis:6379",
		DefaultQueue:  "machinery_tasks",
		ResultBackend: "redis://redis:6379",
		Redis:         &config.RedisConfig{},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		log.Printf("Could not create server: %s", err.Error())
		os.Exit(1)
	}

	if registerErr := server.RegisterTask("vm_deploy", VMDeploy); registerErr != nil {
		log.Printf("Could not register task: %s", registerErr.Error())
		os.Exit(1)
	}

	worker := server.NewWorker("worker-1", 5)
	err = worker.Launch()
	if err != nil {
		log.Printf("Could not launch worker: %s", err.Error())
		os.Exit(1)
	}

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

func VMDeploy(params string) error {
	sDec, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return errors.Wrap(err, "could not decode parameters from base64")
	}

	r := bytes.NewReader(sDec)
	dec := gob.NewDecoder(r)

	var deployParams VMDeployRequest
	err = dec.Decode(&deployParams)
	if err != nil {
		return errors.Wrap(err, "could not decode parameters from bytes")
	}

	return nil
}
