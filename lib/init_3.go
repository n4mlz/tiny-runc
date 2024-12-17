package lib

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func Init_3(containerID string) {
	container, err := LoadContainer(containerID)
	if err != nil {
		log.Fatalf("failed to load container: %v", err)
	}

	config, err := ParseConfig(filepath.Join(container.State.Bundle, "config.json"))
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	Uts(config)

	cmd := exec.Command(config.Process.Args[0], config.Process.Args[1:]...)

	cmd.Env = config.Process.Env
	cmd.Dir = config.Process.Cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: config.Process.User.UID,
			Gid: config.Process.User.GID,
		},
	}

	// TODO: wait for start signal
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run shell: %v", err)
	}
}
