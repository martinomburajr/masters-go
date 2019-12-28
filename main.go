package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const SimulationFilePath = "./_simulation/simulation.json"

func main() {
	paramsPtr := flag.String("params", "", "Pass in the file path (.json) for the given parameters")

	parallelismPtr := flag.Bool("parallelism", true, "Set to false to disable parallelism")
	flag.Parse()

	if *paramsPtr == "" {
		log.Fatal("Params path cannot be empty")
	}

	paramsFolder := *paramsPtr
	parallelism := *parallelismPtr
	log.Println("PARAMS Folder: " + paramsFolder)
	log.Printf("Parallelism Enabled: %t\n", *parallelismPtr)

	//SPEW(paramsFolder, 10)
	Scheduler(paramsFolder, parallelism)
}



// SPEW is used to create the various param files
func SPEW(paramsFolder string, split int) {
	s := simulation.Simulation{}
	abs, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s.SpewJSON(abs, paramsFolder, split)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// SetArguments performs the setup of the simulation and param files
func SetArguments(simulationFilePath, paramsFilePath, dataPath string) (simulation.Simulation,
	evolution.EvolutionParams, error) {
	// Parse
	simulationFile, err := os.Open(simulationFilePath)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	paramsFile, err := os.Open(paramsFilePath)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	var sim simulation.Simulation
	var params evolution.EvolutionParams

	err = json.NewDecoder(simulationFile).Decode(&sim)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	sim.RPath = fmt.Sprintf("%s%s", absolutePath, "/R/runScript.R")
	sim.DataPath = dataPath

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	return sim, params, nil
}

// scheduler runs the actual simulation
func Scheduler(paramsFolder string, parallelism bool) {
	// 1. Get files in the _params3 folder

	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}

	paramFiles, paramFilePath, dataFiles := getParamFiles(absolutePath, paramsFolder)

	filepath.Walk(fmt.Sprintf("%s/%s", absolutePath, "data"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else {
				if info.IsDir() {
					dF := strings.Replace(path, absolutePath+"/data/", "", -1)
					dataFiles = append(dataFiles, dF)
				}
			}
			return nil
		})

	//paramFiles = paramFiles[1:]
	dataFiles = dataFiles[1:]
	//paramFilePath = paramFilePath[1:]

	doneChan := make(chan string, len(paramFiles))
	errChan := make(chan error)

	if parallelism {
		wg := sync.WaitGroup{}
		for i, paramFile := range paramFiles {
			wg.Add(1)
			go func(i int, paramFile string, group *sync.WaitGroup) {
				defer group.Done()
				runSimulation(paramFile, dataFiles, absolutePath, errChan, paramFilePath, i, parallelism, doneChan)
			}(i, paramFile, &wg)
		}
		wg.Done()
	} else {
		for i, paramFile := range paramFiles {
			runSimulation(paramFile, dataFiles, absolutePath, errChan, paramFilePath, i, parallelism, doneChan)
		}
	}

	log.Println("WAIT GROUP COMPLETE")

	for g := range errChan {
		log.Println(g.Error())
	}
}

func runSimulation(paramFile string, dataFiles []string, absolutePath string, errChan chan error, paramFilePath []string, i int, parallelism bool, doneChan chan string) {
	if !contains(paramFile, dataFiles) {
		// create folder and add the started flag
		// add the complete flag
		dataDir := fmt.Sprintf("%s/data/%s", absolutePath, paramFile)
		err := os.MkdirAll(dataDir, 0775)
		if err != nil {
			errChan <- err
		}

		startedPath := fmt.Sprintf("%s/data/%s/%s", absolutePath, paramFile, "started.txt")

		file, err := os.Create(startedPath)
		if err != nil {
			errChan <- err
			log.Println(err.Error())
		}
		_, err = file.WriteString(time.Now().Format(time.RFC3339))
		if err != nil {
			errChan <- err
			log.Println(err.Error())
		}

		simulationn, params, err := SetArguments(SimulationFilePath, paramFilePath[i], dataDir)
		if err != nil {
			errChan <- err
			log.Println(err)
			return
		}
		params.EnableParallelism = parallelism

		newParams, err := simulationn.Begin(params)
		if err != nil {
			errChan <- err
		}

		err, errChan = writeParamFile(absolutePath, paramFile, err, errChan, newParams)

		// add the complete flag
		completePath := fmt.Sprintf("%s/data/%s/%s", absolutePath, paramFile, "completed.txt")
		file2, err := os.Create(completePath)
		if err != nil {
			errChan <- err
			log.Println(err.Error())
		}
		_, err = file2.WriteString(time.Now().Format(time.RFC3339))
		if err != nil {
			errChan <- err
			log.Println(err.Error())
		}

		m := fmt.Sprintf("%s\n%s", "Parsing Complete", paramFile)
		log.Println(m)
		doneChan <- m
	} else {
		m := fmt.Sprintf("%s\n%s", "SKIPPING", paramFile)
		log.Println(m)
		doneChan <- m
	}
}

func writeParamFile(absolutePath string, paramFile string, err error, errChan chan error, params evolution.EvolutionParams) (error, chan error) {
	paramsDataPath := fmt.Sprintf("%s/data/%s/%s", absolutePath, paramFile, "_params.json")
	paramsFile, err := os.Create(paramsDataPath)
	if err != nil {
		errChan <- err
		log.Println(err.Error())
	}
	// pass the params file
	err = json.NewEncoder(paramsFile).Encode(&params)
	if err != nil {
		errChan <- err
		log.Println(err.Error())
	}
	return err, errChan
}

func getParamFiles(absolutePath string, paramsFolder string) ([]string, []string, []string) {
	var paramFiles []string
	var paramFilePath []string
	var dataFiles []string
	filepath.Walk(fmt.Sprintf("%s/%s", absolutePath, paramsFolder),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else {
				if info.IsDir() {

				} else {
					dF := strings.Replace(path, absolutePath+"/"+paramsFolder+"/", "", -1)
					dF = strings.Replace(dF, ".json", "", -1)
					paramFiles = append(paramFiles, dF)
					paramFilePath = append(paramFilePath, path)
				}
			}
			return nil
		})
	return paramFiles, paramFilePath, dataFiles
}

func contains(str string, arr []string) bool {
	for i := range arr {
		if str == arr[i] {
			return true
		}
	}
	return false
}

