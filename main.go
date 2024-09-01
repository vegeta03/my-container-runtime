package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codeclysm/extract/v4"
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
	progress := make(chan float64)
	go func() {
	    for p := range progress {
	        fmt.Printf("Extraction progress: %.2f%%\r", p)
	    }
	}()
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

func extractTar(src string, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	
	ctx := context.Background()
	return extract.Tar(ctx, file, dst, nil)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}