package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseMemoryKvFile(t *testing.T) {
	cgDir = "./"
	service := "s"
	controllerFile := "memory.stat"
	serviceStats := make(map[string]float64)

	err := parseMemoryKvFile(service, controllerFile, serviceStats)
	if err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted memory.stat file")
	}

	err = os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	f, err := os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}

	_, err = f.WriteString("anon 354068")
	if err != nil {
		t.Error(err)
	}

	err = parseMemoryKvFile(service, controllerFile, serviceStats)
	if err != nil {
		t.Error(err)
	}
	if serviceStats["anon"] != 354068 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 354068", serviceStats["anon"])
	}

	controllerFile = "memory.events"
	f, err = os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}

	_, err = f.WriteString("oom 1")
	if err != nil {
		t.Error(err)
	}
	err = parseMemoryKvFile(service, controllerFile, serviceStats)
	if err != nil {
		t.Error(err)
	}
	if serviceStats["oom"] != 1 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 1", serviceStats["oom"])
	}
}

func TestParseMemoryFile(t *testing.T) {
	cgDir = "./"
	service := "s"
	testFile := "memory.low"
	serviceStats := make(map[string]float64)

	err := os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	f, err := os.Create(fmt.Sprint(service, "/", testFile))
	if err != nil {
		t.Error(err)
	}
	_, err = f.WriteString("354068")
	if err != nil {
		t.Error(err)
	}

	err = parseMemoryFile(service, testFile, serviceStats)
	if err != nil {
		t.Error(err)
	}

	if serviceStats["memory.low"] != 354068 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 354068", serviceStats["memory.low"])
	}
}
