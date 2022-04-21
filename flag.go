package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var (
		d time.Duration
		f float64
	)
	flag.DurationVar(&d, "dur", 1*time.Second, "duration flag")
	flag.Float64Var(&f, "float", 0.1, "float flag")
	flag.Parse()
	fmt.Println(d, f)
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags

func main() {
	flag.Var(&myFlags, "list1", "Some description for this param.")
	flag.Parse()
}
