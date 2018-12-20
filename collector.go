package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	// Register memory metrics with prometheus
	prometheus.MustRegister(memoryAnon)
	prometheus.MustRegister(memoryFile)
	prometheus.MustRegister(memoryKernelStack)
	prometheus.MustRegister(memorySlab)
	prometheus.MustRegister(memorySock)
	prometheus.MustRegister(memoryShmem)
	prometheus.MustRegister(memoryFileMapped)
	prometheus.MustRegister(memoryFileDirty)
	prometheus.MustRegister(memoryFileWriteback)
	prometheus.MustRegister(memoryInactiveAnon)
	prometheus.MustRegister(memoryActiveAnon)
	prometheus.MustRegister(memoryInactiveFile)
	prometheus.MustRegister(memoryActiveFile)
	prometheus.MustRegister(memoryUnevictable)
	prometheus.MustRegister(memorySlabReclaimable)
	prometheus.MustRegister(memorySlabUnreclaimable)
	prometheus.MustRegister(memoryPgfault)
	prometheus.MustRegister(memoryPgmajfault)
	prometheus.MustRegister(memoryPgrefill)
	prometheus.MustRegister(memoryPgscan)
	prometheus.MustRegister(memoryPgsteal)
	prometheus.MustRegister(memoryPgactivate)
	prometheus.MustRegister(memoryPgdeactivate)
	prometheus.MustRegister(memoryPglazyfree)
	prometheus.MustRegister(memoryPglazyfreed)
	prometheus.MustRegister(memoryWorkingsetRefault)
	prometheus.MustRegister(memoryWorkingsetActivate)
	prometheus.MustRegister(memoryWorkingsetNodereclaim)
	prometheus.MustRegister(memoryCurrent)
	prometheus.MustRegister(memoryHigh)
	prometheus.MustRegister(memoryLow)
	prometheus.MustRegister(memoryMax)
	prometheus.MustRegister(memoryMin)

	// Register cpu metrics with prometheus
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuNrPeriods)
	prometheus.MustRegister(cpuNrThrottled)
	prometheus.MustRegister(cpuThrottled)
}

func cgroupMetics(hasMemoryController bool, hasCPUController bool) {
	for {
		cgItems := cgServices()

		for _, item := range cgItems {
			if hasMemoryController {
				go cgroupMemoryMetics(item)
			}

			if hasCPUController {
				go cgroupCPUMetics(item)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
