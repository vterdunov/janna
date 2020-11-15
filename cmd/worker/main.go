package main

import (
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"

	w "github.com/vterdunov/janna/internal/worker"
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

	vmDeployFunc := w.VMDeploy{}

	if regErr := server.RegisterTask("vm_deploy", vmDeployFunc.Execute); regErr != nil {
		log.Printf("Could not register task: %s", regErr.Error())
		os.Exit(1)
	}

	vmListFunc := w.VMList{}

	if regErr := server.RegisterTask("vm_list", vmListFunc.Execute); regErr != nil {
		log.Printf("Could not register task: %s", regErr.Error())
		os.Exit(1)
	}

	vmInfoFunc := w.VMDeploy{}

	if regErr := server.RegisterTask("vm_info", vmInfoFunc.Execute); regErr != nil {
		log.Printf("Could not register task: %s", regErr.Error())
		os.Exit(1)
	}

	worker := server.NewWorker("worker-1", 5)
	err = worker.Launch()
	if err != nil {
		log.Printf("Could not launch worker: %s", err.Error())
		os.Exit(1)
	}

}
