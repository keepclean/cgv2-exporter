package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseIOKvFile(service, f string, serviceStats map[string]map[string]float64) error {
	file, err := os.Open(filepath.Join(cgDir, service, "io.stat"))
	if err != nil {
		return err
	}
	defer closeFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 5 {
			return errors.New("Invalid io.stat file format")
		}

		device := devices[fields[0]]
		for _, substring := range fields[1:] {
			if serviceStats[device] == nil {
				serviceStats[device] = map[string]float64{}
			}
			kv := strings.Split(substring, "=")
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				v = float64(0)
			}
			serviceStats[device][kv[0]] = v
		}
	}

	return scanner.Err()
}
