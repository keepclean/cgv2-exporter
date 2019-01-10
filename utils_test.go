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

func TestHasController(t *testing.T) {
	cgroupFile, err := os.Create("cgroup.subtree_control")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(cgroupFile.Name())

	_, err = cgroupFile.WriteString("memory io gpu cpu ram")
	if err != nil {
		t.Error(err)
	}

	cgDir = "./system.slice"
	func() {
		_, err := hasController("memory")
		if err == nil {
			t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted file")
		}
	}()

	cgDir = "./"
	func() {
		_, err := hasController("gpu")
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestSystemdServices(t *testing.T) {
	cgDir = "./folder"
	func() {
		_, err := systemdServices()
		if err == nil {
			t.Error("Something goes wrong because it’s impossible to receive useful information from unexisted folder")
		}
	}()

	cgDir = "./"
	func() {
		_, err := systemdServices()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestControllerFiles(t *testing.T) {
	serviceName := "sssss"
	err := os.MkdirAll(serviceName, 0755)
	if err != nil {
		t.Error(err)
	}

	controllerFile := "io.stat"
	_, err = os.Create(fmt.Sprint(serviceName, "/", controllerFile))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(serviceName)

	cgDir = "./folder"
	func() {
		_, err := controllerFiles("io", serviceName)
		if err == nil {
			t.Error("Something goes wrong because it’s impossible to recieve useful information from unexisted folder")
		}
	}()

	cgDir = "./"
	func() {
		files, err := controllerFiles("io", serviceName)
		if err != nil {
			t.Error(err)
		}
		if files[0] != controllerFile {
			t.Errorf("Getting elements for %q controller is failed", "io")
		}
	}()
}
