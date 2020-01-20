package main

import (
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/analysis"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	repeatDelayPtr := flag.Int64("repeatDelay", 45, "Number of minutes to wait on a file that has already been set")
	spewPtr := flag.Int64("spew", 0, "Creates the set of parameter files, if the value is less than 1, "+
		"it will not spew")
	folderPtr := flag.Int64("folder", 0, "Folder")
	completedStatsPtr := flag.Bool("showProgress", false, "Shows the progress of completed/unstarted/incomplete files")
	stealPtr := flag.Bool("steal", true, "Should steal completed files and automatically back them up")
	coalesceBestPath := flag.String("coalesceBest", "", "Feed in the _dataBackup directory to create a coalescedBest." +
		"csv")

	flag.Parse()

	if *coalesceBestPath != "" && len(*coalesceBestPath) > 3 {
		finalCSV, err := analysis.ReadCSVFile(*coalesceBestPath)
		if err != nil {
			log.Fatal(err)
		}

		outputFilePath := fmt.Sprintf("%s/%s", *coalesceBestPath, "coalescedBest.csv")
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()
		err = gocsv.MarshalFile(finalCSV, outputFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	//fmt.Printf(csvBestAll)
	time.Sleep(time.Second * 10)

	if *paramsPtr == "" {
		log.Fatal("Params path cannot be empty")
	}

	completedStats := *completedStatsPtr
	paramsFolder := *paramsPtr
	parallelism := *parallelismPtr
	folder := *folderPtr
	spew := *spewPtr
	logging := *loggingPtr
	runStats := *runStatsPtr
	dataDir := *dataPtr
	workers := *workerPtr
	repeatDelay := *repeatDelayPtr
	steal := *stealPtr


	abs, _ := filepath.Abs(".")
	if completedStats {

		ShowProgress(abs, paramsFolder, dataDir, repeatDelay)
		return
	}

	log.Println("Parameter Folder: " + paramsFolder)
	log.Println("Data Folder: " + dataDir)
	log.Println("Spew Count: " + strconv.FormatInt(spew, 10))
	log.Println("Worker Count: " + strconv.FormatInt(workers, 10))
	log.Println("Repeat Delay: " + strconv.FormatInt(repeatDelay, 10))
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

	Scheduler(paramsFolder, dataDir, parallelism, workers, repeatDelay,steal, logging, runStats)
}

func ShowProgress(abs string, paramsFolder string, dataDir string, repeatDelay int64) {
	completeParamFolder, unstartedParams, incompleteParams :=
		GetParamFileStatus(abs, paramsFolder, dataDir, repeatDelay)
	fileCount := len(completeParamFolder) + len(unstartedParams) + len(incompleteParams)
	msg := fmt.Sprintf("\n%s ==>\n\t\tNumber of Complete Simulations: \t\t (%d/%d)\n"+
		"\t\tNumber of Incomplete Simulations: \t\t(%d/%d)\n"+
		"\t\tNumber of Unstarted Simulations: \t\t (%d/%d)\n",
		time.Now().Format(time.RFC850),
		len(completeParamFolder),
		fileCount,
		len(incompleteParams),
		fileCount,
		len(unstartedParams),
		fileCount)
	log.Printf(msg)
}

func StealCompleted(abs string, paramsFolder string, dataDir, backupFolder, backupParams string, repeatDelay int64) {
	backupDataPath := fmt.Sprintf("%s/%s", abs, backupFolder)
	backupParamsPath := fmt.Sprintf("%s/%s", abs, backupParams)
	os.Mkdir(backupDataPath, 0775)
	os.Mkdir(backupParamsPath, 0775)

	for {
		completeParamFolder, _, _ :=
			GetParamFileStatus(abs, paramsFolder, dataDir, repeatDelay)
			if len(completeParamFolder) < 1 {
				continue
			}

		for _, complete := range completeParamFolder {
			newDataBackupPath := fmt.Sprintf("%s/%s", backupDataPath, complete)
			newParamBackupPath := fmt.Sprintf("%s/%s.json", backupParamsPath, complete)
			oldParamPath := fmt.Sprintf("%s/%s/%s.json", abs,paramsFolder, complete)
			oldDataPath := fmt.Sprintf("%s/%s/%s", abs,dataDir, complete)

			err := os.MkdirAll(newDataBackupPath, 0775)
			if err != nil {
				log.Println(err)
			}
			gsplit := strings.Split(newParamBackupPath, "/")
			str1 := strings.Join(gsplit[:len(gsplit)-1], "/")

			err = os.MkdirAll(str1, 0775)
			if err != nil {
				log.Println(err)
			}

			mut := sync.Mutex{}
			mut.Lock()
			err = filepath.Walk(oldDataPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					oldDataContent, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}

					dataFileNewPath := fmt.Sprintf("%s/%s", newDataBackupPath, info.Name())
					err = ioutil.WriteFile(dataFileNewPath, oldDataContent, 0755)
					if err != nil {
						return err
					}
				}
				return err
			})
			split := strings.Split(oldDataPath, "/")
			parent := split[:len(split)-1]
			parent2 := strings.Join(parent, "/")
			err = os.RemoveAll(parent2)

			// Copy Param Files
			err = filepath.Walk(oldParamPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					oldParamData, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}

					err = ioutil.WriteFile(newParamBackupPath, oldParamData, 0755)
					if err != nil {
						return err
					}
				}
				return err
			})
			splitParam := strings.Split(oldParamPath, "/")
			parentParam := strings.Join(splitParam[:len(splitParam)-1], "/")
			err = os.RemoveAll(parentParam)
			mut.Unlock()
		}
		time.Sleep(time.Second * 2)
	}
}