// +build linux

package main

import "golang.org/x/sys/unix"

func totalRAMMemory() uint64 {
	info := &unix.Sysinfo_t{}
	unix.Sysinfo(info)
	return info.Totalram
}
