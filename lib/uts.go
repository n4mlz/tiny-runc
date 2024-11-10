package lib

import "syscall"

func setHostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}

func setDomainname(domainname string) error {
	return syscall.Setdomainname([]byte(domainname))
}
