package main

import (
	"fmt"
	"os"
)

func new_file() {
	path := "open_file_by_go"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

func main() {
	fmt.Println("hello")
}
