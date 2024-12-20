package lib

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Init_1(containerID, pipeFromParent string, pipeToParent string) {
	container, err := LoadContainer(containerID)
	if err != nil {
		fmt.Println("Error loading container:", err)
		return
	}

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

	cmd := exec.Command("/proc/self/exe", "init", "2", container.State.ID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
