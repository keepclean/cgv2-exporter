package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func parseKV(s string) (string, uint64, error) {
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return "", 0, errors.New("Invalid format")
	}

	v, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return fields[0], v, nil

}

func closeFile(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func cgServices() (items []string) {
	entries, err := ioutil.ReadDir(cgDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range entries {
		if !item.IsDir() {
			continue
		}
		if !strings.HasSuffix(item.Name(), ".service") {
			continue
		}
		items = append(items, item.Name())
	}

	return
}

func hasController(c string) bool {
	file, err := ioutil.ReadFile(filepath.Join(cgDir, "cgroup.subtree_control"))
	if err != nil {
		log.Fatalf("Can't check availability cgroups controllers: %v", err)
	}

	return strings.Contains(string(file), c)
}

func totalRAMMemory() uint64 {
	info := &unix.Sysinfo_t{}
	unix.Sysinfo(info)
	return info.Totalram
}

func controllerFiles(controller, item string) ([]string, error) {
	entries, err := ioutil.ReadDir(filepath.Join(cgDir, item))
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), fmt.Sprint(controller, ".")) {
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}
