// +build linux

package main

import "golang.org/x/sys/unix"

func totalRAMMemory() float64 {
	info := &unix.Sysinfo_t{}
	unix.Sysinfo(info)
	return float64(info.Totalram)
}
