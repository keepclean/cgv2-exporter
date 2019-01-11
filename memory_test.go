package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMemoryKvFile(t *testing.T) {
	cgDir = "./"
	service := "s"
	controllerFile := "memory.stat"
	stat := &memoryStat{}

	err := parseMemoryKvFile(service, controllerFile, stat)
	if err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted io.stat file")
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

	err = parseMemoryKvFile(service, controllerFile, stat)
	if err != nil {
		t.Error(err)
	}
	if stat.Anon != 354068 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 354068", stat.Anon)
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
	err = parseMemoryKvFile(service, controllerFile, stat)
	if err != nil {
		t.Error(err)
	}
	if stat.EventsOom != 1 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 1", stat.Anon)
	}
}

func TestParseMemoryFiles(t *testing.T) {
	cgDir = "./"
	service := "s"
	testFiles := []string{"memory.stat", "memory.low"}

	err := os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	for _, testFile := range testFiles {
		f, err := os.Create(fmt.Sprint(service, "/", testFile))
		if err != nil {
			t.Error(err)
		}

		if testFile == "memory.stat" {
			_, err = f.WriteString("file 354068")
		}
		if testFile == "memory.low" {
			_, err = f.WriteString("354068")
		}
		if err != nil {
			t.Error(err)
		}
	}

	stat := &memoryStat{}
	parseMemoryFiles(service, stat)
	if stat.Low != 354068 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 354068", stat.Low)
	}
	if stat.File != 354068 {
		t.Errorf("Something wrong with parsing test memory.stat file: got %v, want 354068", stat.File)
	}
}

func TestCgroupMemoryMetrics(t *testing.T) {
	cgDir = "./"
	service := "s"
	testFiles := []string{"memory.stat", "memory.max"}

	err := os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)
	for _, testFile := range testFiles {
		f, err := os.Create(fmt.Sprint(service, "/", testFile))
		if err != nil {
			t.Error(err)
		}

		if testFile == "memory.stat" {
			_, err = f.WriteString("slab 1")
		}
		if testFile == "memory.max" {
			_, err = f.WriteString("2")
		}
		if err != nil {
			t.Error(err)
		}
	}

	// TODO Make test a bit smarter
	cgroupMemoryMetrics(service, false)
	if _, err := memoryMax.GetMetricWithLabelValues(service); err != nil {
		t.Error("Something goes wrong... ", err)
	}

	cgroupMemoryMetrics(service, true)
	if _, err := memoryMax.GetMetricWithLabelValues(service); err != nil {
		t.Error("Something goes wrong... ", err)
	}
}
