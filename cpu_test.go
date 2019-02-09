package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseCPUKvFile(t *testing.T) {
	cgDir = "./"
	service := "s"
	file := "cpu.stat"
	stats := make(map[string]float64)

	err := parseCPUKvFile(service, file, stats)
	if err == nil {
		t.Errorf("Something goes wrong because it’s impossible to receive useful information from unexisted %s file", file)
	}

	err = os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Create(fmt.Sprint(service, "/", file))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = f.WriteString("usage_usec 55002\nnr_periods 1")
	if err != nil {
		t.Error(err)
	}

	if err = parseCPUKvFile(service, file, stats); err != nil {
		t.Error(err)
	}

	if stats["usage_usec"] != (float64(55002) / 1e6) {
		t.Errorf("Something wrong with parsing test cpu.stat file: got %v, want %f", stats["usage_usec"], (float64(55002) / 1e6))
	}
	if stats["nr_periods"] != 1 {
		t.Errorf("Something wrong with parsing test cpu.stat file: got %v, want 1", stats["nr_periods"])
	}
}
