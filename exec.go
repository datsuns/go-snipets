package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Printf(scanner.Text())
		}
		done <- true
	}()
	cmd.Start()
	<-done
	err = cmd.Wait()
}
