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
	flag.Parse()

	if *paramsPtr == "" {
		log.Fatal("Params path cannot be empty")
	}

	paramsFolder := *paramsPtr
	log.Println("PARAMS Folder: " + paramsFolder)

	//SPEW(paramsFolder)

	scheduler(paramsFolder)
}

func SPEW(paramsFolder string) {
	s := simulation.Simulation{}
	abs, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s.SpewJSON(abs, paramsFolder, 6)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ParseInputArguments allows the user to pass in the simulation and evolution parameters into the system to begin
// processing.
func ParseInputArguments() (simulation.Simulation, evolution.EvolutionParams, error) {
	simulationPtr := flag.String("simulation", "", "Pass in the file path (.json) for the given simulation")
	paramsPtr := flag.String("params", "", "Pass in the file path (.json) for the given parameters")
	flag.Parse()

	if *simulationPtr == "" {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf("simulation .json file must be specified")
	}
	if *paramsPtr == "" {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf("parameter .json file must be specified")
	}

	// Parse
	simulationFile, err := os.Open(*simulationPtr)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	paramsFile, err := os.Open(*paramsPtr)
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

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	paramsPtr = nil
	simulationPtr = nil

	return sim, params, nil
}

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

func scheduler(paramsFolder string) {
	// 1. Get files in the _params3 folder
	var paramFiles []string
	var paramFilePath []string
	var dataFiles []string

	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}

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

	paramFiles = paramFiles[1:]
	dataFiles = dataFiles[1:]
	paramFilePath = paramFilePath[1:]

	doneChan := make(chan string, len(paramFiles))
	errChan := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(len(paramFiles))
	for i, paramFile := range paramFiles {
		go func(i int, paramFile string, group *sync.WaitGroup) {
			defer group.Done()
			if !contains(paramFile, dataFiles) {
				// create folder and add the started flag
				// add the complete flag
				dataDir := fmt.Sprintf("%s/data/%s",absolutePath,paramFile)
				err := os.MkdirAll(dataDir, 0775)
				if err != nil {
					errChan <- err
				}

				startedPath := fmt.Sprintf("%s/data/%s/%s",absolutePath,paramFile,"started.txt")

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

				_, err = simulationn.Begin(params)
				if err != nil {
					errChan <- err
					//log.Fatal(err)
				}

				// add the complete flag
 				completePath := fmt.Sprintf("%s/data/%s/%s",absolutePath,paramFile,"completed.txt")
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

		}(i, paramFile, &wg)
	}
	wg.Done()

	log.Println("WAIT GROUP COMPLETE")

	for g := range errChan {
		log.Println(g.Error())
	}
}

func contains(str string, arr []string) bool {
	for i := range arr {
		if str == arr[i] {
			return true
		}
	}
	return false
}


func folderContainsFile(file, folder string) bool {
	contains := false
	filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else {
				if strings.Contains(path, file) {
					contains = true
				}
			}
			return nil
		})
	return contains
}