package main

import (
	"flag"
	"fmt"
	"github.com/martinomburajr/masters-go/analysis"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const SimulationFilePath = "./_simulation/simulation.json"

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Set seed

	paramsPtr := flag.String("params", "_params", "Pass in the file path (.json) for the given parameters")
	dataPtr := flag.String("dataDir", "data", "Pass in the file path (.json) for the given parameters")
	parallelismPtr := flag.Bool("parallelism", true, "Set to false to disable parallelism")
	loggingPtr := flag.Bool("logging", true, "Should Log to stdout and logs.logs file")
	runStatsPtr := flag.Bool("runStats", true, "Can run R based statistics")
	workerPtr := flag.Int64("numWorkers", 2, "Number of workers (each attaches to a paramfile)")
	repeatDelayPtr := flag.Int64("repeatDelay", 45, "Number of minutes to wait on a file that has already been set")
	spewPtr := flag.Int64("spew", -1, "Creates the set of parameter files, if the value is less than 1, "+
		"it will not spew")
	folderPtr := flag.Int64("folder", 0, "Folder")
	completedStatsPtr := flag.Bool("showProgress", false, "Shows the progress of completed/unstarted/incomplete files")
	stealPtr := flag.Bool("steal", true, "Should steal completed files and automatically back them up")
	rIndependentParentDir := flag.String("runRIndependent", "", "run's are to a given set of directories. "+
		"The value supplied must be the parent folder containing all the folders that require R to run in.")
	//--analysisBaseFolder="/home/martinomburajr/Desktop/Results"
	analyisBaseFolder := flag.String("analysisBaseFolder", "", "pass the base folder containing all the different simulations. This will coalesce relevant files")
	runFolder := flag.String("runFolder", "", "pass in the paramFolder to run, " +
		"do not pass in the parent folder e.g. TopologySET-4")
	flag.Parse()


	if  *rIndependentParentDir != "" {
		dirs := make([]string, 0)
		finalDires := make([]string, 0)
		err := filepath.Walk(*rIndependentParentDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return err
		})
		dirs = dirs[1:]
		for i := 1; i < len(dirs); i+=2 {
			finalDires = append(finalDires, dirs[i])
		}
		if err != nil {
			log.Fatalf("RunR: %s", err.Error())
		}
		RunR(finalDires)
	}

	//
	if *analyisBaseFolder != "" {
		wg := sync.WaitGroup{}
		wg.Add(4)
		//errChan := make(chan error)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := analysis.CombineGenerations(*analyisBaseFolder)
			if err != nil {
				fmt.Println(err.Error())
				//errChan <- err
			}
		}(&wg)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := analysis.CombineBest(*analyisBaseFolder)
			if err != nil {
				fmt.Println(err.Error())
				//errChan <- err
			}
		}(&wg)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
		err := analysis.CombineBestCombinedAll(*analyisBaseFolder)
		if err != nil {
			fmt.Println(err.Error())
			//errChan <- err
		}
		}(&wg)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
		err := analysis.CombineBestCombinedAll2(*analyisBaseFolder)
		if err != nil {
			fmt.Println(err.Error())
			//errChan <- err
		}
		}(&wg)

		wg.Wait()
		fmt.Println("F")
		return
	}

	//fmt.Printf(csvBestAll)

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

	os.Mkdir(paramsFolder, 0777)

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
	log.Printf("Run Folder: %s\n", *runFolder)

	if spew > 0 {
		SPEWNoSplit(paramsFolder)
		return
	}
	if *runFolder != "" {
		err := SimpleScheduler(*runFolder, dataDir, logging, runStats)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	Scheduler(paramsFolder, dataDir, parallelism, workers, repeatDelay, steal, logging, runStats)
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
			oldParamPath := fmt.Sprintf("%s/%s/%s.json", abs, paramsFolder, complete)
			oldDataPath := fmt.Sprintf("%s/%s/%s", abs, dataDir, complete)

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
			filepath.Walk(oldParamPath, func(path string, info os.FileInfo, err error) error {
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
			os.RemoveAll(parentParam)
			mut.Unlock()
		}
		time.Sleep(time.Second * 2)
	}
}

func steal(paramToStealFolder, abs string, paramsFolder string, dataDir, backupFolder, backupParams string,
) error {
	backupDataPath := fmt.Sprintf("%s/%s", abs, backupFolder)
	backupParamsPath := fmt.Sprintf("%s/%s", abs, backupParams)
	os.Mkdir(backupDataPath, 0775)
	os.Mkdir(backupParamsPath, 0775)

	newDataBackupPath := fmt.Sprintf("%s/%s", backupDataPath, paramToStealFolder)
	newParamBackupPath := fmt.Sprintf("%s/%s.json", backupParamsPath, paramToStealFolder)
	oldParamPath := fmt.Sprintf("%s/%s/%s.json", abs, paramsFolder, paramToStealFolder)
	oldDataPath := fmt.Sprintf("%s/%s/%s", abs, dataDir, paramToStealFolder)

	mut := sync.Mutex{}
	mut.Lock()
	err := filepath.Walk(oldDataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			oldDataContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			dataFileNewPath := fmt.Sprintf("%s/%s", newDataBackupPath, info.Name())
			os.MkdirAll(newDataBackupPath, 0777)

			err = ioutil.WriteFile(dataFileNewPath, oldDataContent, 0777)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	split := strings.Split(oldDataPath, "/")
	parent := split[:len(split)-1]
	parent2 := strings.Join(parent, "/")
	err = os.RemoveAll(parent2)
	if err != nil {
		return err
	}

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
			newParamBackupFolder := strings.ReplaceAll(newParamBackupPath, ".json", "")
			os.MkdirAll(newParamBackupFolder, 0777)
			err = ioutil.WriteFile(newParamBackupPath, oldParamData, 0777)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	splitParam := strings.Split(oldParamPath, "/")
	parentParam := strings.Join(splitParam[:len(splitParam)-1], "/")
	err = os.RemoveAll(parentParam)
	if err != nil {
		return err
	}
	mut.Unlock()

	return nil
}