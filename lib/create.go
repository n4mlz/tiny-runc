package lib

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Create() {
	// TODO: get path from args according to OCI Runtime Specification
	c := ParseConfig("bundle/config.json")

	cmd := exec.Command("/proc/self/exe", "init")

	var cloneFlags uintptr
	for _, ns := range c.Linux.Namespaces {
		if ns.Type == "pid" {
			cloneFlags |= syscall.CLONE_NEWPID
		}
		if ns.Type == "ipc" {
			cloneFlags |= syscall.CLONE_NEWIPC
		}
		if ns.Type == "uts" {
			cloneFlags |= syscall.CLONE_NEWUTS
		}
		if ns.Type == "mount" {
			cloneFlags |= syscall.CLONE_NEWNS
		}
		if ns.Type == "user" {
			cloneFlags |= syscall.CLONE_NEWUSER
		}
		// TODO: add more namespaces (network, cgroup, time)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: cloneFlags,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run command: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
