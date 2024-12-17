package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func Create(containerID string, bundlePath string) {
	if bundlePath == "" {
		bundlePath = "."
	}

	container, err := NewContainer(containerID, bundlePath)
	if err != nil {
		fmt.Println("Error creating container:", err)
		return
	}

	config, err := ParseConfig(filepath.Join(container.State.Bundle, "config.json"))
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	// TODO: properly manage pipe path
	pipeToChild := "tmp/parent_to_child"
	pipeFromChild := "tmp/child_to_parent"

	if err := SetupPipes(pipeToChild, pipeFromChild); err != nil {
		fmt.Println("Error setting up pipes:", err)
		return
	}
	defer CleanupPipes(pipeToChild, pipeFromChild)

	cmd := exec.Command("/proc/self/exe", "init", "1", container.State.ID, pipeToChild, pipeFromChild)

	var cloneFlags uintptr
	for _, ns := range config.Linux.Namespaces {
		if ns.Type == "user" {
			cloneFlags |= syscall.CLONE_NEWUSER
		}
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: cloneFlags,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	handler, err := OpenPipesA(pipeFromChild, pipeToChild)
	if err != nil {
		fmt.Println("Error opening pipes:", err)
		return
	}
	defer handler.Close()

	message, err := handler.ReceiveMessage()
	if err != nil {
		panic(err)
	} else if message != "ready" {
		panic("child not ready")
	}

	// TODO: newuidmap, newgidmap
	// TODO: read from config
	uidMapPath := fmt.Sprintf("/proc/%d/uid_map", cmd.Process.Pid)
	gidSetGroupPath := fmt.Sprintf("/proc/%d/setgroups", cmd.Process.Pid)
	gidMapPath := fmt.Sprintf("/proc/%d/gid_map", cmd.Process.Pid)

	if err := os.WriteFile(uidMapPath, []byte("0 1000 1"), 0600); err != nil {
		panic(err)
	}
	if err := os.WriteFile(gidSetGroupPath, []byte("deny"), 0600); err != nil {
		panic(err)
	}
	if err := os.WriteFile(gidMapPath, []byte("0 1000 1"), 0600); err != nil {
		panic(err)
	}

	if err = handler.SendMessage("done"); err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	cmd.Wait()
}
