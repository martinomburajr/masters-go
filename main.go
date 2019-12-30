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
	"runtime"
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

	fmt.Printf("Current Goroutines: %d", runtime.NumGoroutine())

	paramsFolder := *paramsPtr
	parallelism := *parallelismPtr
	log.Println("PARAMS Folder: " + paramsFolder)
	log.Printf("Parallelism Enabled: %t\n", *parallelismPtr)

	//SPEW(paramsFolder, 4)
	parallelism = true
	Scheduler(paramsFolder, parallelism, true, true)
}


// scheduler runs the actual simulation
func Scheduler(paramsFolder string, parallelism, logging, runStats bool) {
	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}

	paramFiles := getParamFiles(absolutePath, paramsFolder)
	dataFiles := getDataFiles(absolutePath)
	dataFiles = dataFiles[1:]

	sim := simulationParams{
		absolutePath:  absolutePath,
		dataFiles:     dataFiles,
		paramFolder: paramsFolder,
		errChan:       make(chan error),
		logChan:       make(chan string),
		doneChan:       make(chan bool),
		parallelism:   parallelism,
		logging:   logging,
		runStats: runStats,
	}

	// Listen to logs and errors

	go func(simulationParam *simulationParams) {
		started := time.Now()
		for {
			select {
			case logg := <- simulationParam.logChan:
				fmt.Println("LOG: " + logg)
			case err := <- simulationParam.errChan:
				fmt.Println("Error: " + err.Error())
				return
			case isDone := <- simulationParam.doneChan:
				elapsedTime := time.Since(started)

				msg := fmt.Sprintf("\nElapsed Time: %s\nIsComplete: %t\n", elapsedTime.String(), isDone)
				fmt.Println(msg)

				close(sim.doneChan)
				os.Exit(0)
			}
		}
	}(&sim)


	if parallelism {
		wg := sync.WaitGroup{}
		for _, paramFile := range paramFiles {
			wg.Add(1)

			go func(sim simulationParams, paramFile string, group *sync.WaitGroup) {
				defer group.Done()
				sim.paramFile = paramFile
				runSimulation(sim)
			}(sim, paramFile, &wg)
		}
		wg.Wait()
	} else {
		for _, paramFile := range paramFiles {
			sim.paramFile = paramFile
			runSimulation(sim)
		}
	}

	close(sim.logChan)
	close(sim.errChan)
}

type simulationParams struct {
	absolutePath string
	paramFile    string
	paramFolder  string
	dataFiles    []string
	errChan      chan error
	logChan      chan string

	parallelism bool
	logging     bool
	runStats    bool
	doneChan    chan bool
}

func runSimulation(simulationParams simulationParams) {
	paramFilePath := fmt.Sprintf("%s/%s/%s.json", simulationParams.absolutePath,
		simulationParams.paramFolder, simulationParams.paramFile)

	if !contains(simulationParams.paramFile, simulationParams.dataFiles) {
		dataDir := fmt.Sprintf("%s/data/%s", simulationParams.absolutePath, simulationParams.paramFile)
		err := os.MkdirAll(dataDir, 0775)
		if err != nil {
			simulationParams.errChan <- err
		}

		// started
		createFileInDataDir(simulationParams, "started.txt", time.Now().Format(time.RFC3339))

		simulationn, params, err := SetArguments(SimulationFilePath, paramFilePath, dataDir)
		if err != nil {
			if simulationParams.parallelism {
				simulationParams.errChan <- err
			}
			log.Printf(err.Error())
			return
		}

		params.EnableParallelism = simulationParams.parallelism
		params.EnableLogging = simulationParams.logging
		params.RunStats = simulationParams.runStats
		params.LoggingChan = simulationParams.logChan
		params.ErrorChan = simulationParams.errChan
		params.DoneChan = simulationParams.doneChan

		newParams, err := simulationn.Begin(params)
		if err != nil {
			if simulationParams.parallelism {
				simulationParams.errChan <- err
			}
			log.Printf(err.Error())
		}

		writeParamFile(simulationParams, newParams, simulationParams.errChan)

		// completed
		createFileInDataDir(simulationParams, "completed.txt", time.Now().Format(time.RFC3339))

		m := fmt.Sprintf("%s\n%s", "Parsing Complete", simulationParams.paramFile)
		if simulationParams.parallelism {
			simulationParams.logChan <- m
		}else {
			log.Printf(m)
		}

	} else {
		m := fmt.Sprintf("%s\n%s", "SKIPPING", simulationParams.paramFile)
		if simulationParams.parallelism {
			ch := len(simulationParams.logChan)
			fmt.Println(ch)
			if simulationParams.logChan != nil {
				simulationParams.logChan <- m
			}
		} else {
			log.Printf(m)
		}

	}
}

func createFileInDataDir(simulationParams simulationParams, filename, content string) {
	completePath := fmt.Sprintf("%s/data/%s/%s", simulationParams.absolutePath, simulationParams.paramFile, filename)

	file, err := os.Create(completePath)
	if err != nil {
		simulationParams.errChan <- err
	}
	_, err = file.WriteString(content)
	if err != nil {
		simulationParams.errChan <- err
	}
}

func writeParamFile(sim simulationParams, params evolution.EvolutionParams, errChan chan error) {
	paramsDataPath := fmt.Sprintf("%s/data/%s/%s", sim.absolutePath, sim.paramFile, "_params.json")
	paramsFile, err := os.Create(paramsDataPath)

	if err != nil {
		errChan <- err
	}

	// pass the params file
	err = json.NewEncoder(paramsFile).Encode(&params)
	if err != nil {
		errChan <- err
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

	sim.RPath = fmt.Sprintf("%s%s", absolutePath, "/R")
	sim.DataPath = dataPath

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	return sim, params, nil
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


func getParamFiles(absolutePath string, paramsFolder string) (paramFiles []string) {
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
				}
			}
			return nil
		})
	return paramFiles
}

func getDataFiles(absolutePath string) (dataFiles []string) {
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

	return dataFiles
}

func contains(str string, arr []string) bool {
	for i := range arr {
		if str == arr[i] {
			return true
		}
	}
	return false
}

