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

	uidMappingArgs := []string{fmt.Sprint(cmd.Process.Pid)}
	for _, mapping := range config.Linux.UIDMappings {
		uidMappingArgs = append(uidMappingArgs, fmt.Sprint(mapping.ContainerID), fmt.Sprint(mapping.HostID), fmt.Sprint(mapping.Size))
	}

	gidMappingArgs := []string{fmt.Sprint(cmd.Process.Pid)}
	for _, mapping := range config.Linux.GIDMappings {
		gidMappingArgs = append(gidMappingArgs, fmt.Sprint(mapping.ContainerID), fmt.Sprint(mapping.HostID), fmt.Sprint(mapping.Size))
	}

	if err := exec.Command("newuidmap", uidMappingArgs...).Run(); err != nil {
		panic(err)
	}

	if err := exec.Command("newgidmap", gidMappingArgs...).Run(); err != nil {
		panic(err)
	}

	if err = handler.SendMessage("done"); err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	cmd.Wait()
}
