package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/vterdunov/janna/internal/virtualmachine"

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

	if err := server.RegisterTask("vm_deploy", VMDeploy); err != nil {
		log.Printf("Could not register task: %s", err.Error())
		os.Exit(1)
	}

	worker := server.NewWorker("worker-1", 5)
	err = worker.Launch()
	if err != nil {
		log.Printf("Could not launch worker: %s", err.Error())
		os.Exit(1)
	}

}

func VMDeploy(params string) error {
	sDec, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return errors.Wrap(err, "could not decode parameters from base64")
	}

	r := bytes.NewReader(sDec)
	dec := gob.NewDecoder(r)

	var deployParams virtualmachine.VMDeployRequest
	err = dec.Decode(&deployParams)
	if err != nil {
		return errors.Wrap(err, "could not decode parameters from bytes")
	}

	return nil
}

// func sliceToMap(slice []string) map[string]string {
// 	sliceLen := len(slice)
// 	resMap := make(map[string]string, sliceLen)

// 	for i := 0; i < len(slice); i += 2 {
// 		if i == sliceLen-1 {
// 			break
// 		}

// 		resMap[slice[i]] = slice[i+1]
// 	}

// 	return resMap
// }
