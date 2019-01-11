// +build darwin

package main

// Fake func for running "go test" on macos
func totalRAMMemory() uint64 {
	return 99999999
}
