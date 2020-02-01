package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func RunR(pathToDir []string) {
	RFiles := []string{
		"best.R", "best-combined.R", "epochs.R",
		"generations.R", "strategy.R",
	}

	doneCounterChan := make(chan bool, len(pathToDir))
	for i := range pathToDir {
		wg := sync.WaitGroup{}
		for _, rFile := range RFiles {
			wg.Add(1)
			go func(group *sync.WaitGroup, rFile string, dirPath string) {
				defer group.Done()

				RPath := "R"

				fqdn := fmt.Sprintf("%s/%s", RPath, rFile)
				cmd := exec.Command("Rscript", fqdn, dirPath)
				msg := fmt.Sprintf("Rscript: \n%s\n", cmd.String())
				fmt.Println(msg)
				//logChan <- msg

				err := cmd.Run()
				if err != nil {
					fmt.Printf("error: %s", err.Error())
				}
			}(&wg, rFile, pathToDir[i])
		}
		wg.Wait()
		doneCounterChan <- true
		fmt.Printf("\n \t %s", pathToDir[i])
	}

	counter := 0
	for {
		select {
		case <-doneCounterChan:
			counter++
			fmt.Printf("MAIN: Completed (%d/%d) files", counter, len(pathToDir))

			if counter == len(pathToDir)-1 {
				break
			}
		}
		if counter == len(pathToDir)-1 {
			break
		}
	}
	fmt.Println("COMPLETED R Independent")
}
