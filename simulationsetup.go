package main

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// scheduler runs the actual simulation
func Scheduler(paramsFolder, dataDirName string, parallelism bool, numberOfSimultaneousParams, repeatDelay int64,
	logging,
	runStats bool) {
	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}

	dataFiles := getDataFiles(absolutePath, dataDirName)
	dataFiles = dataFiles[1:]

	sim := simulationParams{
		absolutePath:               absolutePath,
		dataFiles:                  dataFiles,
		dataDirName:                dataDirName,
		numberOfSimultaneousParams: int(numberOfSimultaneousParams),
		paramFolder:                paramsFolder,
		errChan:                    make(chan error),
		logChan:                    make(chan string),
		doneChan:                   make(chan bool),
		parallelism:                parallelism,
		logging:                    logging,
		runStats:                   runStats,
	}

	// Listen to logs and errors
	go SetupLogger(sim)

	completeParamFolder, unstartedParams, incompleteParams :=
		GetParamFileStatus(sim.absolutePath, sim.paramFolder, sim.dataDirName, repeatDelay)
	fileCount := len(completeParamFolder) + len(unstartedParams) + len(incompleteParams)

	msg := fmt.Sprintf("\nNumber of Complete Simulations: (%d/%d)\n"+
		"Number of Incomplete Simulations: (%d/%d)\n"+
		"Number of Unstarted Simulations: (%d/%d)\n",
		len(completeParamFolder),
		fileCount,
		len(incompleteParams),
		fileCount,
		len(unstartedParams),
		fileCount)
	log.Printf(msg)

	time.Sleep(3 * time.Second)
	if len(unstartedParams) == 0 && len(incompleteParams) == 0 {
		log.Printf("\n\n################################### NO WORK TO DO! ###################################\n\n")
		return
	} else {
		count := 0
		for len(completeParamFolder) <= fileCount {
			//combinedunfinished := append(unstartedParams, incompleteParams...)
			if len(unstartedParams) <= sim.numberOfSimultaneousParams {
				sim.numberOfSimultaneousParams = len(unstartedParams)
			}
			if len(unstartedParams) != 0 && sim.numberOfSimultaneousParams == 0 {
				fmt.Printf("\n\nSTARTING NEW WORK\n\n")
				if len(unstartedParams) > int(numberOfSimultaneousParams) {
					sim.numberOfSimultaneousParams =  int(numberOfSimultaneousParams)
				} else {
					sim.numberOfSimultaneousParams = len(unstartedParams)
				}
			}
			if len(unstartedParams) == 0 {
				fmt.Printf("\n\n################################### WAITING ON" +
					" OTHERS TO HOPEFULLY COMPLETE! ###################################\n\n")
			}
			WorkerPool(sim.numberOfSimultaneousParams, len(unstartedParams),
				func(index int, waitGroup *sync.WaitGroup) {
					if index >= len(unstartedParams) {
						mut := sync.Mutex{}
						mut.Lock()
						index = rand.Intn(len(unstartedParams))
						mut.Unlock()
					}
					runSimulation(sim, unstartedParams[index])
				})

			completeParamFolder, unstartedParams, incompleteParams = GetParamFileStatus(sim.absolutePath,
				sim.paramFolder, sim.dataDirName, repeatDelay)
			log.Printf("\n\n\n################################### COMPLETED CYCLE %d"+
				"! ###################################\n\n\n", count)
			time.Sleep(5 * time.Second)
			count++
		}
	}

	sim.doneChan <- true
	close(sim.logChan)
	close(sim.errChan)
}

func SetupLogger(simulationParam simulationParams) {
	file := setupLogFile(simulationParam)
	started := time.Now()
	doneCounter := 0
	for {
		select {
		case logg := <-simulationParam.logChan:
			sb := strings.Builder{}
			sb.WriteString("\n" + time.Now().Format(time.RFC3339))
			sb.WriteString(" ==> Folder: ")
			sb.WriteString(simulationParam.paramFile)
			sb.WriteString(" | ")
			sb.WriteString(logg)
			loggg := sb.String()

			fmt.Fprintf(file, loggg)
			fmt.Fprintf(os.Stdout, loggg)
		case err := <-simulationParam.errChan:
			fmt.Println("Error: " + err.Error())
			return
		case isDone := <-simulationParam.doneChan:
			if isDone {
				doneCounter--
			}else {
				doneCounter++
			}
			fmt.Printf("DONE COUNTER: %d", doneCounter)
			if doneCounter <= 0 {
				elapsedTime := time.Since(started)
				msg := fmt.Sprintf("\nElapsed Time: %s\nIsComplete: %t\n", elapsedTime.String(), isDone)
				fmt.Println(msg)
				close(simulationParam.doneChan)

				fmt.Println("GRACEFULLY STAYING UP FOR 2 HOURS")
				time.Sleep(2 * time.Hour)
				os.Exit(0)
			}
		}
	}
}

type simulationParams struct {
	absolutePath               string
	paramFile                  string
	paramFolder                string
	dataDirName                string
	dataFiles                  []string
	errChan                    chan error
	logChan                    chan string
	numberOfSimultaneousParams int

	parallelism bool
	logging     bool
	runStats    bool
	doneChan    chan bool
}

func runSimulation(simulationParams simulationParams, paramFileToRun string) {
	simulationParams.paramFile = paramFileToRun
	simulationParams.doneChan <- false

	dataDir := fmt.Sprintf("%s/data/%s", simulationParams.absolutePath, simulationParams.paramFile)
	err := os.MkdirAll(dataDir, 0775)

	// started
	createFileInDataDir(simulationParams, "started.txt", time.Now().Format(time.RFC3339))
	paramFilePath := fmt.Sprintf("%s/%s/%s.json", simulationParams.absolutePath,
		simulationParams.paramFolder, simulationParams.paramFile)

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
	params.ParamFile = simulationParams.paramFile
	params.StatisticsOutput.OutputPath = simulationn.DataPath

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
	return
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
	params.StatisticsOutput.OutputPath = sim.DataPath

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	return sim, params, nil
}

func setupLogFile(simulationParam simulationParams) *os.File {
	logFolder := "logs"
	logFilePath := fmt.Sprintf("%s/folder-%s-logs.txt", logFolder, simulationParam.paramFile)
	os.Mkdir(logFolder, 0775)
	file, err := os.Create(logFilePath)
	if err != nil {
		simulationParam.errChan <- err
	}
	return file
}
