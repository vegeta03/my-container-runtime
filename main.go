package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		child()
	default:
		fmt.Printf("Unknown command: %[1]s\n", command)
		os.Exit(1)
	}
}

func run() {
	fmt.Printf("Running command: %[1]v\n", os.Args[2:])

	/*
		`/proc/self/exe`: This is a special file in Linux systems that represents the currently running executable.
		By using this as the command to execute, the program is essentially launching a new instance of itself.
	*/
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up namespaces for the child process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("From child, Running command: %[1]v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cg()

	// Set up the container environment
	must(syscall.Sethostname([]byte("container")))
	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Append a string to the pwd
	ubi9Rootfs := filepath.Join(pwd, "/ubi9.4-rootfs/")
	must(syscall.Chroot(ubi9Rootfs))
	must(os.Chdir("/"))

	/*
		The proc filesystem (procfs) is a special filesystem in Unix-like operating systems
		that presents process information as files in a hierarchical file structure. It provides
		an interface to kernel data structures, allowing processes to be examined and manipulated.
		When you mount the proc filesystem, it creates a directory structure under /proc that
		contains information about running processes, system memory, mounted devices, hardware
		configuration, and other system information. Each running process has its own directory
		under /proc, named after its process ID (PID).
	*/
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""))

	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
	must(syscall.Unmount("tmp", 0))
}

func cg() {
	cgroup := "/sys/fs/cgroup"
	containerDir := filepath.Join(cgroup, "container")

	// Create the container directory if it doesn't exist
	if err := os.MkdirAll(containerDir, 0755); err != nil {
		fmt.Printf("Error creating cgroup directory: %v\n", err)
		return
	}

	files := map[string]string{
		"pids.max": "20",
		// "notify_on_release": "1",
		"cgroup.procs": strconv.Itoa(os.Getpid()),
	}

	for filename, content := range files {
		path := filepath.Join(containerDir, filename)
		if err := os.WriteFile(path, []byte(content), 0700); err != nil {
			fmt.Printf("Error writing to %s: %v\n", path, err)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
