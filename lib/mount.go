package lib

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func mount_fs(container Container, config specs.Spec) {
	for _, m := range config.Mounts {
		flags := 0
		if m.Options != nil {
			for _, o := range m.Options {
				switch o {
				case "bind":
					flags |= syscall.MS_BIND
				case "rbind":
					flags |= syscall.MS_BIND | syscall.MS_REC
				case "ro":
					flags |= syscall.MS_RDONLY
				case "rro":
					flags |= syscall.MS_RDONLY | syscall.MS_REC
				case "rw":
					flags &^= syscall.MS_RDONLY
				case "nosuid":
					flags |= syscall.MS_NOSUID
				case "nodev":
					flags |= syscall.MS_NODEV
				case "noexec":
					flags |= syscall.MS_NOEXEC
				case "remount":
					flags |= syscall.MS_REMOUNT
				case "mand":
					flags |= syscall.MS_MANDLOCK
				case "dirsync":
					flags |= syscall.MS_DIRSYNC
				case "noatime":
					flags |= syscall.MS_NOATIME
				case "nodiratime":
					flags |= syscall.MS_NODIRATIME
				case "relatime":
					flags |= syscall.MS_RELATIME
				case "strictatime":
					flags |= syscall.MS_STRICTATIME
				case "silent":
					flags |= syscall.MS_SILENT
				case "posixacl":
					flags |= syscall.MS_POSIXACL
				case "unbindable":
					flags |= syscall.MS_UNBINDABLE
				case "private":
					flags |= syscall.MS_PRIVATE
				case "slave":
					flags |= syscall.MS_SLAVE
				case "shared":
					flags |= syscall.MS_SHARED
				case "rprivate":
					flags |= syscall.MS_PRIVATE | syscall.MS_REC
				case "rslave":
					flags |= syscall.MS_SLAVE | syscall.MS_REC
				case "rshared":
					flags |= syscall.MS_SHARED | syscall.MS_REC
				default:
					log.Printf("unknown mount option %s", o)
				}
			}
		}

		dest := filepath.Join(container.Root, "rootfs", m.Destination)

		if _, err := os.Stat(dest); err != nil {
			if err := os.MkdirAll(dest, 0755); err != nil {
				log.Fatalf("mkdir %s failed: %v", dest, err)
			}
		}

		if err := syscall.Mount(m.Source, dest, m.Type, uintptr(flags), ""); err != nil {
			log.Fatalf("mount %s failed: %v", dest, err)
		}
	}
}

func Mount(container Container, config specs.Spec) {
	// create rootfs path
	rootfs := filepath.Join(container.Root, "rootfs")

	// recursively copy rootfs
	if err := CopyDirectory(filepath.Join(container.State.Bundle, config.Root.Path), rootfs); err != nil {
		log.Fatalf("copy %s failed: %v", config.Root.Path, err)
	}

	// mount("none", "/", NULL, MS_REC | MS_PRIVATE, NULL);
	if err := syscall.Mount("none", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		log.Fatalf("mount none failed: %v", err)
	}

	// mount rootfs
	if err := syscall.Mount(rootfs, rootfs, "", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("mount %s failed: %v", rootfs, err)
	}

	mount_fs(container, config)

	// mkdir old root
	oldRoot := filepath.Join(rootfs, ".oldroot")
	if err := os.MkdirAll(oldRoot, 0755); err != nil {
		log.Fatalf("mkdir %s failed: %v", oldRoot, err)
	}

	// pivot root
	if err := syscall.PivotRoot(rootfs, oldRoot); err != nil {
		log.Fatalf("pivot_root %s %s failed: %v", rootfs, oldRoot, err)
	}

	// chdir to /
	if err := syscall.Chdir("/"); err != nil {
		log.Fatalf("chdir / failed: %v", err)
	}

	// umount old root
	if err := syscall.Unmount("/.oldroot", syscall.MNT_DETACH); err != nil {
		log.Fatalf("umount /.oldroot failed: %v", err)
	}

	// remove old root
	if err := os.Remove("/.oldroot"); err != nil {
		log.Fatalf("remove /.oldroot failed: %v", err)
	}
}
