package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseIOKvFile(service, f string, serviceStats map[string]map[string]float64) error {
	file, err := os.Open(filepath.Join(cgDir, service, f))
	if err != nil {
		return err
	}
	defer closeFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		device := devices[fields[0]]
		for _, substring := range fields[1:] {
			if serviceStats[device] == nil {
				serviceStats[device] = map[string]float64{}
			}
			kv := strings.Split(substring, "=")
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				v = 0
			}
			serviceStats[device][kv[0]] = v
		}
	}

	return scanner.Err()
}
