package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseIOStat(t *testing.T) {
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

	controllerFile := "io.stat"
	f, err := os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = f.WriteString("8:0 rbytes=290260631552 wbytes=0 rios=6740371 wios=354068")
	if err != nil {
		t.Error(err)
	}

	devices = map[string]string{
		"8:0": "sda",
	}

	ioStats, err := parseIOStat(service)
	if err != nil {
		t.Error(err)
	}
	_, ok := ioStats["sda"]
	if !ok {
		t.Errorf("There's no %q key in returned map", "sda")
	}
	v, ok := ioStats["sda"]["rbytes"]
	if !ok {
		t.Errorf("There's no %q key in returned map for %q disk", "rbytes", "sda")
	}
	targetValue := float64(290260631552)
	if v != targetValue {
		t.Errorf("Problems with parsing io.stat file: got %f; wanted: %f", v, targetValue)
	}
}

func TestCgroupIOMetrics(t *testing.T) {
	cgDir = "./"
	service := "s"

	err := os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}

	controllerFile := "io.stat"
	f, err := os.Create(fmt.Sprint(service, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = f.WriteString("8:0 rbytes=290260631552 wbytes=0 rios=6740371 wios=354068")
	if err != nil {
		t.Error(err)
	}

	devices = map[string]string{
		"8:0": "sda",
	}

	// TODO Make test a bit smarter
	cgroupIOMetrics(service, false)
	if _, err := ioRbytes.GetMetricWithLabelValues(service, "sda"); err != nil {
		t.Error("Something goes wrong... ", err)
	}

	cgroupIOMetrics(service, true)
	if _, err := ioRios.GetMetricWithLabelValues(service, "/dev/sda"); err != nil {
		t.Error("Something goes wrong... ", err)
	}

}
