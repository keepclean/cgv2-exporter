package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func parseCPUKvFile(service, f string, serviceStats map[string]float64) error {
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

		if strings.HasSuffix(key, "_usec") {
			serviceStats[key] = value / 1e9
			continue
		}

		serviceStats[key] = value
	}

	return scanner.Err()
}
