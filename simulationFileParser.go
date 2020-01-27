package main

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// GetParamFileStatus returns files that have been thoroughly processed,
// versus those that have not or are in an intermediary state.
func GetParamFileStatus(absolutePath, paramDirName, dataDirName string, repeatDelay int64) (completeParamFolder,
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

	mut := sync.Mutex{}
	mut.Lock()
	if len(dataFiles) > 0 {
		dataFiles = dataFiles[1:]
	}

	for _, dataFile := range dataFiles {
		// Check if at least has reached 5% of generations
		has25Txt := strings.Contains(dataFile, "1.txt")
		if has25Txt {
			split := strings.Split(dataFile, "/")
			str := strings.Builder{}
			str.WriteString(split[0])
			str.WriteString("/")
			str.WriteString(split[1])
			finalString := fmt.Sprintf("%s/%s/%s/%s", absolutePath, dataDirName, str.String(), "1.txt")

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
			if seconds > float64(repeatDelay*60) {
				dataPath2 := fmt.Sprintf("%s/%s/%s", absolutePath, dataDirName, dataPath)
				os.RemoveAll(dataPath2)
				paramDataMap[dataPath] = -1
			} else {
				paramDataMap[dataPath] = 25
			}
		}
	}

	for _, dataFile := range dataFiles {
		hasCompletedTxt := strings.Contains(dataFile, "completed.txt")
		hasFinalParamsFile := strings.Contains(dataFile, "_params.json")
		if hasCompletedTxt || hasFinalParamsFile {
			split := strings.Split(dataFile, "/")
			str := strings.Builder{}
			str.WriteString(split[0])
			str.WriteString("/")
			str.WriteString(split[1])
			finalString := str.String()

			paramDataMap[finalString] = 1
		}
	}

	keys := make([]string, 0)
	for k, _ := range paramDataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if paramDataMap[k] == 1 {
			completeParamFolder = append(completeParamFolder, k)
		} else if paramDataMap[k] == -1 {
			unstardedParamFolder = append(unstardedParamFolder, k)
		} else if paramDataMap[k] == 25 {
			incompleteParamFolder = append(incompleteParamFolder, k)
		}
	}

	mut.Unlock()

	//for _, incompleteFolder := range incompleteParamFolder {
	//	os.RemoveAll(incompleteFolder)
	//}

	return completeParamFolder, unstardedParamFolder, incompleteParamFolder
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

func createFileInDataDir(simulationParams simulationParams, filename, content string) {
	completePath := fmt.Sprintf("%s/data/%s/%s", simulationParams.absolutePath, simulationParams.paramFile, filename)

	mut := sync.Mutex{}
	mut.Lock()
	file, err := os.Create(completePath)
	if err != nil {
		simulationParams.errChan <- err
	}
	_, err = file.WriteString(content)
	if err != nil {
		simulationParams.errChan <- err
	}
	file.Close()
	mut.Unlock()
}

func writeParamFile(sim simulationParams, params evolution.EvolutionParams, errChan chan error) {
	paramsDataPath := fmt.Sprintf("%s/data/%s/%s", sim.absolutePath, sim.paramFile, "_params.json")
	mut := sync.Mutex{}
	mut.Lock()

	paramsFile, err := os.Create(paramsDataPath)
	if err != nil {
		errChan <- err
	}

	// pass the params file
	err = json.NewEncoder(paramsFile).Encode(&params)
	if err != nil {
		errChan <- err
	}
	paramsFile.Close()
	mut.Unlock()
}
