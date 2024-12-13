package lib

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Init_2(containerID string) {
	// TODO: get path from args according to OCI Runtime Specification
	c := ParseConfig("bundle/config.json")

	cmd := exec.Command("/proc/self/exe", "init", "3", containerID)

	var cloneFlags uintptr
	for _, ns := range c.Linux.Namespaces {
		if ns.Type == "pid" {
			cloneFlags |= syscall.CLONE_NEWPID
		}
		if ns.Type == "mount" {
			cloneFlags |= syscall.CLONE_NEWNS
		}
		if ns.Type == "ipc" {
			cloneFlags |= syscall.CLONE_NEWIPC
		}
		if ns.Type == "uts" {
			cloneFlags |= syscall.CLONE_NEWUTS
		}
		if ns.Type == "time" {
			cloneFlags |= syscall.CLONE_NEWTIME
		}
		// if ns.Type == "network" {
		// 	cloneFlags |= syscall.CLONE_NEWNET
		// }

		// if ns.Type == "cgroup" {
		// 	cloneFlags |= syscall.CLONE_NEWCGROUP
		// }
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: cloneFlags,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run shell: %v", err)
	}
}
