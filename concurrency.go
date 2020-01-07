package main

import (
	"log"
	"runtime"
	"sync"
)

func WorkerPool(numWorkers, numIterations int, workerFunction func(index int, waitGroup *sync.WaitGroup)) {
	if numWorkers == 0 {
		return
	}
	if numIterations == 0 {
		return
	}
	numIterations = numIterations + 1
	numWorkers = numWorkers+1
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
		//if didBreak {
		//	break
		//}
	}
}