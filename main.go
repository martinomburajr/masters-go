package main

import (
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"os/exec"
)

func main() {
	//name := "run.json"

	//
	//f, err := os.Create("test.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//json.NewEncoder(f).Encode(params)

	simulation := Simulation{
		NumberOfRunsPerState: 1,
		Name:                 "simulation-1",
	}


	params := evolution.EvolutionParams{

	}


	simulation.Begin(params)
	cmd := exec.Command("Rscript", "launch.R")
	log.Fatal(cmd.Run())
}
