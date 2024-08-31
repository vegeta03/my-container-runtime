package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go run <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		run()
	case "child":
		fmt.Println("Child command received")
	default:
		fmt.Printf("Unknown command: %[1]s\n", command)
		os.Exit(1)
	}
}

func run() {
	fmt.Printf("Running command: %[1]s\n", os.Args[2:])

	/* 
	`/proc/self/exe`: This is a special file in Linux systems that represents the currently running executable. 
	By using this as the command to execute, the program is essentially launching a new instance of itself. 
	*/
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up the container environment
	must(syscall.Sethostname([]byte("container")))
	
	extractDir := "ubi9.4_rootfs"
	err := extractTar("ubi9.4_rootfs.tar", extractDir)
	if err != nil {
		panic(err)
	}

	// Use the extracted directory for Chroot
	must(syscall.Chroot(extractDir))

	err = cmd.Run()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
