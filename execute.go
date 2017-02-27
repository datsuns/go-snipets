package main

import (
	"fmt"
	"os/exec"
)

func execute(tool string, params []string, debug bool) {
	if debug {
		fmt.Printf(" >> %v %v\n", tool, params)
	}
	log, err := exec.Command(tool, params...).Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("log: %s\n", log)
}

func run() {
	params := [][]string{
		{"-i", "'/.*MaxRecords=.*/d'", "PAC.ini"},
		{"-i", "-e", "'s/RingBufferMode=false/RingBufferMode=true/'", "PAC.ini"},
	}
	for _, p := range params {
		execute("sed", p)
	}
}
