package lib

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Init_3(containerID string, bundlePath string) {
	if bundlePath == "" {
		bundlePath = "."
	}

	c := ParseConfig(filepath.Join(bundlePath, "config.json"))

	for _, ns := range c.Linux.Namespaces {
		if ns.Type == "uts" {
			if c.Hostname != "" {
				if err := setHostname(c.Hostname); err != nil {
					log.Fatalf("failed to set hostname: %v", err)
				}
			}
		}
	}

	if c.Domainname != "" {
		if err := setDomainname(c.Domainname); err != nil {
			log.Fatalf("failed to set domainname: %v", err)
		}
	}

	cmd := exec.Command(c.Process.Args[0], c.Process.Args[1:]...)
	cmd.Env = c.Process.Env
	cmd.Dir = c.Process.Cwd

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// TODO: wait for start signal
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run shell: %v", err)
	}
}
