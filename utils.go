package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseKV(s string) (string, float64, error) {
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return "", 0, errors.New("Invalid format")
	}

	v, err := strconv.ParseFloat(fields[1], 64)
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

func cgroupFiles(service string) ([]string, error) {
	var files []string
	entries, err := ioutil.ReadDir(filepath.Join(cgDir, service))
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
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
