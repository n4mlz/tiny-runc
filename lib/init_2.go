package lib

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func Init_2(containerID string) {
	container, err := LoadContainer(containerID)
	if err != nil {
		log.Fatalf("failed to load container: %v", err)
	}

	config, err := ParseConfig(filepath.Join(container.State.Bundle, "config.json"))
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	cmd := exec.Command("/proc/self/exe", "init", "3", container.State.ID)

	var cloneFlags uintptr
	for _, ns := range config.Linux.Namespaces {
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
