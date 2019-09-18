package main

import (
	"fmt"
)

var (
	Map map[int]string = map[int]string{
		1:  "hello",
		10: "hi",
	}
)

func main() {
	fmt.Printf("value! [%v]\n", Map[1])
	if v, ok := Map[200]; ok {
		fmt.Println("exits ", v)
	} else {
		fmt.Println("NOT exits")
	}
}
