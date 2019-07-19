package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func main() {
	fmt.Println("Hello Worker!")

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

	server.RegisterTask("add", Add)

	worker := server.NewWorker("worker-1", 5)
	err = worker.Launch()
	if err != nil {
		log.Printf("Could not launch worker: %s", err.Error())
		os.Exit(1)
	}

}

func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}
