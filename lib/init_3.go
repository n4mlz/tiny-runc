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
			if c.Hostname != nil {
				if err := setHostname(*c.Hostname); err != nil {
					log.Fatalf("failed to set hostname: %v", err)
				}
			}
		}
	}

	if c.Domainname != nil {
		if err := setDomainname(*c.Domainname); err != nil {
			log.Fatalf("failed to set domainname: %v", err)
		}
	}

	// TODO: read from config
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run shell: %v", err)
	}
}
