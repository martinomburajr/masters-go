package main

import (
	"github.com/martinomburajr/masters-go/simulation"
	"log"
	"path/filepath"
)

// SPEW is used to create the various param files. Split refers to the number of folders to create
//func SPEW(paramsFolder string, split int) {
//	s := simulation.Simulation{}
//	abs, err := filepath.Abs(".")
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	err = s.SpewJSON(abs, paramsFolder, split)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//}

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
