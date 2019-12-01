// +build linux

package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func totalRAMMemory() float64 {
	info := &unix.Sysinfo_t{}
	err := unix.Sysinfo(info)
	if err != nil {
		log.Println(err)
	}
	return float64(info.Totalram)
}
