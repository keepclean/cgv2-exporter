package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseIOKvFile(t *testing.T) {
	cgDir = "./"
	service := "s"
	file := "io.stat"
	serviceIOStats := make(map[string]map[string]float64)

	if err := parseIOKvFile(service, file, serviceIOStats); err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted io.stat file")
	}

	if err := os.Mkdir(service, 0755); err != nil {
		t.Error(err)
	}

	f, err := os.Create(fmt.Sprint(service, "/", file))
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

	err = parseIOKvFile(service, file, serviceIOStats)
	if err != nil {
		t.Error(err)
	}
	_, ok := serviceIOStats["sda"]
	if !ok {
		t.Errorf("There's no %q key in returned map", "sda")
	}
	v, ok := serviceIOStats["sda"]["rbytes"]
	if !ok {
		t.Errorf("There's no %q key in returned map for %q disk", "rbytes", "sda")
	}
	targetValue := float64(290260631552)
	if v != targetValue {
		t.Errorf("Problems with parsing io.stat file: got %f; wanted: %f", v, targetValue)
	}
}
