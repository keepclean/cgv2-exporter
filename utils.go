package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func systemdServices() ([]string, error) {
	var services []string
	entries, err := ioutil.ReadDir(cgDir)
	if err != nil {
		return services, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".service") {
			continue
		}
		services = append(services, entry.Name())
	}

	return services, nil
}

func cgroupControllers() (controllers map[string]bool) {
	controllers = make(map[string]bool)
	for _, c := range []string{"memory", "cpu", "io"} {
		v, err := hasController(c)
		if err != nil {
			log.Println(err)
			continue
		}
		controllers[c] = v
	}

	return
}

func hasController(c string) (bool, error) {
	file, err := ioutil.ReadFile(filepath.Join(cgDir, "cgroup.subtree_control"))
	if err != nil {
		return false, fmt.Errorf("Can't check availability %q cgroups controllers: %v", c, err)
	}

	return strings.Contains(string(file), c), nil
}

func controllerFiles(controller, service string) ([]string, error) {
	entries, err := ioutil.ReadDir(filepath.Join(cgDir, service))
	if err != nil {
		return []string{}, err
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

var devices map[string]string

func blockDevices() {
	sysBlockDir := "/sys/block/"
	entries, err := ioutil.ReadDir(sysBlockDir)
	if err != nil {
		log.Println(err)
	}

	devices = make(map[string]string)
	for _, entry := range entries {
		if entry.Mode()&os.ModeSymlink == 0 {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(sysBlockDir, entry.Name(), "dev"))
		if err != nil {
			log.Println(err)
			continue
		}

		device := strings.TrimSpace(string(file))
		devices[device] = entry.Name()
	}

	return
}
