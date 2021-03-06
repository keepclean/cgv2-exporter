package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseKV(t *testing.T) {
	// case #1 - test wrong string
	s := "k1 v1 k2 v2"
	_, _, err := parseKV(s)
	if err == nil {
		t.Error("ParseKV is broken in case #1")
	}

	// case #2 - test string with right format but wrong content
	s = "k1 v1"
	_, _, err = parseKV(s)
	if err == nil {
		t.Error("ParseKV is broken in case #2")
	}

	// case #3 - test correct string
	s = "k1 2"
	k, v, err := parseKV(s)
	if err != nil || k != "k1" || v != 2 {
		t.Error("ParseKV is broken in case #3")
	}
}

func TestSystemdServices(t *testing.T) {
	cgDir = "./folder"
	_, err := systemdServices()
	if err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted folder")
	}

	cgDir = "./"
	_, err = systemdServices()
	if err != nil {
		t.Error(err)
	}
}

func TestCgroupFiles(t *testing.T) {
	cgDir = "./folder"
	service := "s"
	_, err := cgroupFiles(service)
	if err == nil {
		t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted folder")
	}

	cgDir = "./"
	err = os.Mkdir(service, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(service)

	_, err = os.Create(fmt.Sprint(service, "/", "memory.stat"))
	if err != nil {
		t.Error(err)
	}

	_, err = cgroupFiles(service)
	if err != nil {
		t.Error(err)
	}
}
