package lib

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Child() {
	// TODO: get path from args according to OCI Runtime Specification
	c := ParseConfig("bundle/config.json")

	pipeToParent := os.Args[2]
	pipeFromParent := os.Args[3]

	handler, err := OpenPipesB(pipeFromParent, pipeToParent)
	if err != nil {
		fmt.Println("Error opening pipes:", err)
		return
	}
	defer handler.Close()

	if err = handler.SendMessage("ready"); err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	message, err := handler.ReceiveMessage()
	if err != nil {
		panic(err)
	} else if message != "done" {
		panic("parent not done")
	}

	syscall.Setuid(0)
	syscall.Setgid(0)

	cmd := exec.Command("/proc/self/exe", "init")

	var cloneFlags uintptr
	for _, ns := range c.Linux.Namespaces {
		if ns.Type == "uts" {
			cloneFlags |= syscall.CLONE_NEWUTS
		}
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: cloneFlags,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
