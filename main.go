package main

import (
	"flag"
	"log"
	"strconv"
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
	repeatDelay := *repeatDelayPtr
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

	Scheduler(paramsFolder, dataDir, parallelism, workers, repeatDelay,logging, runStats)
}
