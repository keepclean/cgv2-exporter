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

	// cadvisor style memory metrics for the backward compability
	prometheus.Register(memoryCache)
	prometheus.Register(memoryFailCnt)
	prometheus.Register(memoryMaxUsage)
	prometheus.Register(memoryUsage) // done
	prometheus.Register(memoryRss)
	prometheus.Register(memorySwap)
	prometheus.Register(memoryWorkingSet)           // done
	prometheus.Register(memorySpecLimit)            // done
	prometheus.Register(memorySpecReservationLimit) // unified cgroup doesnt't have anything related
	prometheus.Register(memorySpecSwapLimit)        // unified cgroup doesnt't have anything related
	prometheus.Register(memoryCadvisorPgfault)      // done
	prometheus.Register(memoryCadvisorPgmajfault)   // done

	// Register cpu metrics with prometheus
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuNrPeriods)
	prometheus.MustRegister(cpuNrThrottled)
	prometheus.MustRegister(cpuThrottled)
}

func cgroupMetrics(hasMemoryController bool, hasCPUController bool, cadvisorMemoryMetrics bool) {
	for {
		cgItems := cgServices()

		for _, item := range cgItems {
			if hasMemoryController {
				go cgroupMemoryMetrics(item, cadvisorMemoryMetrics)
			}

			if hasCPUController {
				go cgroupCPUMetrics(item)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
