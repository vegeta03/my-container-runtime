package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		fmt.Println("Running a container!!!")
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
