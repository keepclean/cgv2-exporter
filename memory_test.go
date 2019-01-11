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
}
