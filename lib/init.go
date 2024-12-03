package lib

import (
	"log"
	"os"
	"os/exec"
)

func Init() {
	// TODO: get path from args according to OCI Runtime Specification
	c := ParseConfig("bundle/config.json")

	if c.Hostname != nil {
		if err := setHostname(*c.Hostname); err != nil {
			log.Fatalf("failed to set hostname: %v", err)
		}
	}

	if c.Domainname != nil {
		if err := setDomainname(*c.Domainname); err != nil {
			log.Fatalf("failed to set domainname: %v", err)
		}
	}

	// TODO: implement

	log.Printf("init")

	// run shell
	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run shell: %v", err)
	}
}
