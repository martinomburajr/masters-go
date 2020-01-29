package main

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/masters-go/evolog"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var SimulationArgs = simulation.Simulation{
	NumberOfRunsPerState: 5,
	Name:                 "",
	StatsFiles: []string{
		"best.R", "best-combined.R", "epochs.R",
		"generations.R", "strategy.R",
	},
	RPath: "/R",
}

// scheduler runs the actual simulation
func Scheduler(paramsFolder, dataDirName string, parallelism bool, numberOfSimultaneousParams, repeatDelay int64,
	canSteal,
	logging,
	runStats bool) {
	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}

	dataFiles := getDataFiles(absolutePath, dataDirName)

	if len(dataFiles) > 0 {
		dataFiles = dataFiles[1:]
	}

	sim := simulationParams{
		absolutePath:               absolutePath,
		dataFiles:                  dataFiles,
		dataDirName:                dataDirName,
		numberOfSimultaneousParams: int(numberOfSimultaneousParams),
		paramFolder:                paramsFolder,
		errChan:                    make(chan error),
		logChan:                    make(chan evolog.Logger),
		doneChan:                   make(chan bool),
		parallelism:                parallelism,
		logging:                    logging,
		runStats:                   runStats,
	}

	// Listen to logs and errors
	go SetupLogger(sim)

	if canSteal {
		go StealCompleted(sim.absolutePath, paramsFolder, sim.dataDirName, "_dataBackup", "_paramsBackup", repeatDelay)
	}

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

	if len(unstartedParams) == 0 && len(incompleteParams) == 0 {
		log.Printf("\n\n################################### NO WORK TO DO! ###################################\n\n")
		return
	} else {
		count := 0
		for len(completeParamFolder) <= fileCount && count < fileCount {
			sim.doneChan <- false
			if len(unstartedParams) != 0 && unstartedParams != nil {
				runSimulation(sim, unstartedParams[0])
			}
			sim.doneChan <- true

			completeParamFolder, unstartedParams, incompleteParams = GetParamFileStatus(sim.absolutePath,
				sim.paramFolder, sim.dataDirName, repeatDelay)
			log.Printf("\n\n\n################################### COMPLETED CYCLE %d"+
				"! ###################################\n\n\n", count)
			count++
		}
	}

	sim.doneChan <- true
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("GRACEFULLY STAYING UP FOR 12 Minutes")
		time.Sleep(12 * time.Minute)
	}(&wg)
	wg.Wait()
	close(sim.logChan)
	close(sim.errChan)
}

func SetupLogger(simulationParam simulationParams) {
	file := setupLogFile(simulationParam)
	defer file.Close()

	started := time.Now()
	doneCounter := 0

	for {
		select {
		case logg := <-simulationParam.logChan:
			sb := strings.Builder{}
			sb.WriteString("\n" + time.Now().Format(time.RFC3339))
			sb.WriteString(simulationParam.paramFile)
			sb.WriteString(" | ")
			sb.WriteString(logg.Message)
			loggg := sb.String()

			fmt.Fprintf(file, loggg)
		case err := <-simulationParam.errChan:
			msg := fmt.Sprintf("Error: %s" + err.Error())
			fmt.Println(msg)
			fmt.Fprintf(file, msg)
			//return
		case isDone, ok := <-simulationParam.doneChan:
			if !ok {
				elapsedTime := time.Since(started)
				msg := fmt.Sprintf("\nElapsed Time: %s\nIsComplete: %t\n", elapsedTime.String(), isDone)
				fmt.Println(msg)
				close(simulationParam.doneChan)

				fmt.Println("GRACEFULLY STAYING UP FOR 12 Minutes")
				time.Sleep(12 * time.Minute)
				os.Exit(0)
			}
			if isDone {
				doneCounter--
			} else {
				doneCounter++
			}
			fmt.Printf("DONE COUNTER: %d\n", doneCounter)
			if doneCounter < 0 {
				elapsedTime := time.Since(started)
				msg := fmt.Sprintf("\nElapsed Time: %s\nIsComplete: %t\n", elapsedTime.String(), isDone)
				fmt.Println(msg)
				close(simulationParam.doneChan)

				fmt.Println("GRACEFULLY STAYING UP FOR 5 hours")
				time.Sleep(5 * time.Hour)
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
	logChan                    chan evolog.Logger
	numberOfSimultaneousParams int

	parallelism bool
	logging     bool
	runStats    bool
	doneChan    chan bool
}

func runSimulation(simulationParams simulationParams, paramFileToRun string) {
	simulationParams.paramFile = paramFileToRun

	dataDir := fmt.Sprintf("%s/data/%s", simulationParams.absolutePath, simulationParams.paramFile)
	err := os.MkdirAll(dataDir, 0775)

	// started
	createFileInDataDir(simulationParams, "started.txt", time.Now().Format(time.RFC3339))
	paramFilePath := fmt.Sprintf("%s/%s/%s.json", simulationParams.absolutePath,
		simulationParams.paramFolder, simulationParams.paramFile)

	params, err := SetArguments(&SimulationArgs, paramFilePath, dataDir)
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
	params.StatisticsOutput.OutputPath = SimulationArgs.DataPath

	newParams, err := SimulationArgs.Begin(params)
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
func SetArguments(simulation *simulation.Simulation, paramsFilePath, dataPath string) (evolution.EvolutionParams,
	error) {
	paramsFile, err := os.Open(paramsFilePath)
	defer paramsFile.Close()
	if err != nil {
		return evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	var params evolution.EvolutionParams
	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	simulation.RPath = fmt.Sprintf("%s%s", absolutePath, "/R")
	simulation.DataPath = dataPath
	params.StatisticsOutput.OutputPath = simulation.DataPath

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	return params, nil
}

func setupLogFile(simulationParam simulationParams) *os.File {
	logFolder := "logs"
	logFilePath := fmt.Sprintf("%s/folder-logs.txt", logFolder)

	os.Mkdir(logFolder, 0775)
	file, err := os.Create(logFilePath)
	if err != nil {
		simulationParam.errChan <- err
	}
	return file
}
