package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go run <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		fmt.Println("Run command received.")
	case "child":
		fmt.Println("Child command received")
	default:
		fmt.Printf("Unknown command: %[1]s\n", command)
		os.Exit(1)
	}
}
