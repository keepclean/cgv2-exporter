package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var totalRAM = totalRAMMemory()

func parseMemoryKvFile(service, f string, serviceStats map[string]float64) error {
	file, err := os.Open(filepath.Join(cgDir, service, f))
	if err != nil {
		return err
	}
	defer closeFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, err := parseKV(scanner.Text())
		if err != nil {
			return err
		}
		serviceStats[key] = value
	}

	return scanner.Err()
}

func parseMemoryFile(service, f string, serviceStats map[string]float64) error {
	file, err := ioutil.ReadFile(filepath.Join(cgDir, service, f))
	if err != nil {
		return err
	}

	if strings.Contains(string(file), "max") {
		serviceStats[f] = totalRAM
		return nil
	}

	v, err := strconv.ParseFloat(strings.TrimSpace(string(file)), 64)
	if err != nil {
		v = 0
	}
	serviceStats[f] = v

	return nil
}
