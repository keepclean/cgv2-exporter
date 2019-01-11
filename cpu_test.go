package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseCPUStat(t *testing.T) {
	cgDir = "./"
	service := "s"
	_, err := parseIOStat(service)
	if err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted io.stat file")
	}

	err = os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}

	controllerFile := "cpu.stat"
	f, err := os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = f.WriteString("usage_usec 55002")
	if err != nil {
		t.Error(err)
	}

	stat := &cpuStat{}
	if err = parseCPUStat(service, stat); err != nil {
		t.Error(err)
	}

	if stat.Usage != (float64(55002) / 1e9) {
		t.Errorf("Something wrong with parsing test cpu.stat file: got %v, want %f", stat.Usage, (float64(55002) / 1e9))
	}
}

func TestCgroupCPUMetrics(t *testing.T) {
	cgDir = "./"
	service := "s"

	err := os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}

	controllerFile := "cpu.stat"
	f, err := os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = f.WriteString("user_usec 1")
	if err != nil {
		t.Error(err)
	}

	// TODO Make test a bit smarter
	cgroupCPUMetrics(service)
	if _, err = cpuUser.GetMetricWithLabelValues(service); err != nil {
		t.Error("Something goes wrong... ", err)
	}
}
