package lib

import (
	"log"
	"os"
	"os/exec"
)

func Init_3(containerID string) {
	// TODO: get path from args according to OCI Runtime Specification
	c := ParseConfig("bundle/config.json")

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
