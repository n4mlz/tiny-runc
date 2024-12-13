package lib

import (
	"log"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func setHostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}

func setDomainname(domainname string) error {
	return syscall.Setdomainname([]byte(domainname))
}

func Uts(c specs.Spec) {
	for _, ns := range c.Linux.Namespaces {
		if ns.Type == "uts" {
			if c.Hostname != "" {
				if err := setHostname(c.Hostname); err != nil {
					log.Fatalf("failed to set hostname: %v", err)
				}
			}
			if c.Domainname != "" {
				if err := setDomainname(c.Domainname); err != nil {
					log.Fatalf("failed to set domainname: %v", err)
				}
			}
		}
	}
}
