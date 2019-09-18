package main

import (
	"fmt"
	"path/filepath"
)

func glob() {
	files, err := filepath.Glob("*")
	if err != nil {
		panic(err)
	}
	for _, e := range files {
		fmt.Printf("[%v]\n", e)
	}
}

func main() {
	glob()
}
