package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const SimulationFilePath = "./_simulation/simulation.json"

func main() {
	paramsPtr := flag.String("params", "", "Pass in the file path (.json) for the given parameters")
	dataPtr := flag.String("dataDir", "data", "Pass in the file path (.json) for the given parameters")
	parallelismPtr := flag.Bool("parallelism", true, "Set to false to disable parallelism")
	loggingPtr := flag.Bool("logging", true, "Should Log to stdout and logs.logs file")
	runStatsPtr := flag.Bool("runstats", true, "Can run R based statistics")
	workerPtr := flag.Int64("numWorkers", 5, "Number of workers (each attaches to a paramfile)")
	spewPtr := flag.Int64("spew", 0, "Creates the set of parameter files, if the value is less than 1, "+
		"it will not spew")
	folderPtr := flag.Int64("folder", 0, "Folder")

	flag.Parse()

	if *paramsPtr == "" {
		log.Fatal("Params path cannot be empty")
	}

	paramsFolder := *paramsPtr
	parallelism := *parallelismPtr
	folder := *folderPtr
	spew := *spewPtr
	logging := *loggingPtr
	runStats := *runStatsPtr
	dataDir := *dataPtr
	workers := *workerPtr
	log.Println("Parameter Folder: " + paramsFolder)
	log.Println("Data Folder: " + dataDir)
	log.Println("Spew Count: " + strconv.FormatInt(spew, 10))
	log.Println("Worker Count: " + strconv.FormatInt(workers, 10))
	log.Println("Folder: " + strconv.FormatInt(folder, 10))
	log.Printf("Parallelism Enabled: %t\n", *parallelismPtr)
	log.Printf("Logging Enabled: %t\n", *loggingPtr)
	log.Printf("RunStats Enabled: %t\n", *runStatsPtr)

	if spew > 0 {
		if spew == 1 {
			SPEWNoSplit(paramsFolder)
		} else {
			SPEW(paramsFolder, int(spew))
		}
		return
	}

	Scheduler(paramsFolder, dataDir, parallelism, workers, logging, runStats)
}

// scheduler runs the actual simulation
func Scheduler(paramsFolder, dataDirName string, parallelism bool, numberOfSimultaneousParams int64, logging,
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

	if parallelism {
		completeParamFolder, unstartedParams, incompleteParams :=
			GetParamFileStatus(sim.absolutePath, sim.paramFolder, sim.dataDirName)
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
				combinedunfinished := append(unstartedParams, incompleteParams...)
				if len(combinedunfinished) < sim.numberOfSimultaneousParams {
					sim.numberOfSimultaneousParams = len(combinedunfinished)
				}
				WorkerPool(sim.numberOfSimultaneousParams, len(unstartedParams),
					func(index int, waitGroup *sync.WaitGroup) {
						if index >= len(combinedunfinished) {
							mut := sync.Mutex{}
							mut.Lock()
							index = rand.Intn(len(combinedunfinished))
							mut.Unlock()
						}
						runSimulation(sim, combinedunfinished[index])
					})

				completeParamFolder, unstartedParams, incompleteParams =
					GetParamFileStatus(sim.absolutePath, sim.paramFolder, sim.dataDirName)
				log.Printf("\n\n\n################################### COMPLETED CYCLE %d"+
					"! ###################################\n\n\n", count)
				time.Sleep(10 * time.Second)
				count++
			}
		}
	}

	sim.doneChan <- true
	close(sim.logChan)
	close(sim.errChan)
}

func SetupLogger(simulationParam simulationParams) {
	file := setupLogFile(simulationParam)
	started := time.Now()
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
			elapsedTime := time.Since(started)

			msg := fmt.Sprintf("\nElapsed Time: %s\nIsComplete: %t\n", elapsedTime.String(), isDone)
			fmt.Println(msg)
			close(simulationParam.doneChan)
			os.Exit(0)
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

	dataDir := fmt.Sprintf("%s/data/%s", simulationParams.absolutePath, simulationParams.paramFile)
	err := os.MkdirAll(dataDir, 0775)
	if err != nil {
		simulationParams.errChan <- err
	}

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

	m := fmt.Sprintf("%s\n%s", "Parsing Complete", simulationParams.paramFile)
	if simulationParams.parallelism {
		simulationParams.logChan <- m
	} else {
		log.Printf(m)
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
	params.StatisticsOutput.OutputPath = sim.DataPath

	err = json.NewDecoder(paramsFile).Decode(&params)
	if err != nil {
		return simulation.Simulation{}, evolution.EvolutionParams{},
			fmt.Errorf(err.Error())
	}

	return sim, params, nil
}

// SPEW is used to create the various param files. Split refers to the number of folders to create
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

func SPEWNoSplit(paramsFolder string) {
	s := simulation.Simulation{}
	abs, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s.SpewJSONNoSplit(abs, paramsFolder)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getParamFiles(absolutePath string, paramsFolder string) (paramFiles []string) {
	paramPath := fmt.Sprintf("%s/%s", absolutePath, paramsFolder)
	filepath.Walk(paramPath,
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

func getDataFiles(absolutePath string, dataDirName string) (dataFiles []string) {
	inputPath := fmt.Sprintf("%s/%s", absolutePath, dataDirName)
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else {
			if info.IsDir() {
				dataPath := fmt.Sprintf("%s/%s/", absolutePath, dataDirName)
				dF := strings.Replace(path, dataPath, "", -1)

				dataFiles = append(dataFiles, dF)
			}
		}
		return nil
	})
	return dataFiles
}

func getAllDataFiles(absolutePath string, dataDirName string) (dataFiles []string) {
	inputPath := fmt.Sprintf("%s/%s", absolutePath, dataDirName)
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else {
			if !info.IsDir() {
				dataPath := fmt.Sprintf("%s/%s/", absolutePath, dataDirName)
				dF := strings.Replace(path, dataPath, "", -1)

				dataFiles = append(dataFiles, dF)
			}
		}
		return nil
	})

	return dataFiles
}

// GetParamFileStatus returns files that have been thoroughly processed,
// versus those that have not or are in an intermediary state.
func GetParamFileStatus(absolutePath, paramDirName, dataDirName string) (completeParamFolder,
	unstardedParamFolder, incompleteParamFolder []string) {

	completeParamFolder = make([]string, 0)
	unstardedParamFolder = make([]string, 0)
	incompleteParamFolder = make([]string, 0)
	dataFiles := getAllDataFiles(absolutePath, dataDirName)
	paramFiles := getParamFiles(absolutePath, paramDirName)
	paramDataMap := map[string]int{}

	for _, paramFile := range paramFiles {
		paramDataMap[paramFile] = -1
	}

	if len(dataFiles) > 0 {
		dataFiles = dataFiles[1:]
	}

	for _, dataFile := range dataFiles {
		// Check if at least has reached 25% of generations
		has25Txt := strings.Contains(dataFile, "25.txt")
		if has25Txt {
			split := strings.Split(dataFile, "/")
			str := strings.Builder{}
			str.WriteString(split[0])
			str.WriteString("/")
			str.WriteString(split[1])
			finalString := fmt.Sprintf("%s/%s/%s/%s", absolutePath, dataDirName, str.String(), "25.txt")

			str2 := strings.Builder{}
			str2.WriteString(split[0])
			str2.WriteString("/")
			str2.WriteString(split[1])
			dataPath := str2.String()

			file, err := os.Open(finalString)
			if err != nil {
				return nil, nil, nil
			}
			timeStr, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, nil, nil
			}
			parsedTime, err := time.Parse(time.RFC3339, string(timeStr))
			if err != nil {
				return nil, nil, nil
			}

			subtractedTime := time.Now().Sub(parsedTime)
			seconds := subtractedTime.Seconds()
			if seconds > 60*60 {
				paramDataMap[dataPath] = -1
			} else {
				paramDataMap[dataPath] = 25
			}
		}
	}

	for _, dataFile := range dataFiles {
		hasCompletedTxt := strings.Contains(dataFile, "completed.txt")
		if hasCompletedTxt {
			split := strings.Split(dataFile, "/")
			str := strings.Builder{}
			str.WriteString(split[0])
			str.WriteString("/")
			str.WriteString(split[1])
			finalString := str.String()

			paramDataMap[finalString] = 1
		}
	}

	for k, v := range paramDataMap {
		if v == 1 {
			completeParamFolder = append(completeParamFolder, k)
		} else if v == -1 {
			unstardedParamFolder = append(unstardedParamFolder, k)
		} else if v == 25 {
			incompleteParamFolder = append(incompleteParamFolder, k)
		}
	}

	return completeParamFolder, unstardedParamFolder, incompleteParamFolder
}

func contains(str string, arr []string) bool {
	for i := range arr {
		if str == arr[i] {
			return true
		}
	}
	return false
}

func WorkerPool(numWorkers, numIterations int, workerFunction func(index int, waitGroup *sync.WaitGroup)) {
	numIterations = numIterations + 1
	numWorkers = numWorkers + 1
	wg := sync.WaitGroup{}
	timesIn := 0
	didBreak := false
	for !didBreak {
		i := 0
		for {
			mod := i % numWorkers
			if mod != 0 {
				log.Printf("Number of Goroutines: %d", runtime.NumGoroutine())
				wg.Add(1)
				timesIn++
				go workerFunction(timesIn, &wg)
			} else {
				if i != 0 {
					wg.Wait()
				}
			}
			if timesIn == numIterations {
				didBreak = true
				break
			}
			i++
		}
		if didBreak {
			break
		}
	}
}

func moveFolders(sim simulationParams, completedPath, destination string, errChan chan error, logChan chan string,
	wg *sync.WaitGroup) {
	defer wg.Done()
	err := os.Mkdir(destination, 0775)
	//if err != nil {
	//	errChan <- err
	//	return
	//}

	movableDirs := make([]string, 0)

	dataFilePath := fmt.Sprintf("%s/data/%s", sim.absolutePath, completedPath)
	err = filepath.Walk(dataFilePath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".png") {
			//dirs := strings.SplitAfterN(path, "/", 2)
			movableDirs = append(movableDirs, "")
		}
		return err
	})
	if err != nil {
		errChan <- err
		return
	}

	paramFolderToMove := ""
	err = os.Rename(paramFolderToMove, destination)
	if err != nil {
		errChan <- err
		return
	}
	msg := fmt.Sprintf("Moved PARAM Folder %s to %s", paramFolderToMove, destination)
	logChan <- msg
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
