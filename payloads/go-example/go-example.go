package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hello %v\n", os.Args[1:])
	} else {
		fmt.Println("Hello World!")
	}
	fmt.Println("This code was written in Go!")
}
