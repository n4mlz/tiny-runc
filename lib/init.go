package lib

import (
	"log"
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
}
