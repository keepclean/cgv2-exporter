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
	prometheus.MustRegister(memoryEventsLow)
	prometheus.MustRegister(memoryEventsHigh)
	prometheus.MustRegister(memoryEventsMax)
	prometheus.MustRegister(memoryEventsOom)
	prometheus.MustRegister(memoryEventsOomKill)

	// cadvisor style memory metrics for the backward compability
	prometheus.Register(memoryCache)                // done
	prometheus.Register(memoryFailCnt)              // done
	prometheus.Register(memoryMaxUsage)             // unified cgroup doesn'y have anything related to this out of the box
	prometheus.Register(memoryUsage)                // done
	prometheus.Register(memoryRss)                  // done
	prometheus.Register(memorySwap)                 // TODO parse memory.swap.current file
	prometheus.Register(memoryWorkingSet)           // done
	prometheus.Register(memorySpecLimit)            // done
	prometheus.Register(memorySpecReservationLimit) // unified cgroup doesnt't have anything related to this
	prometheus.Register(memorySpecSwapLimit)        // unified cgroup doesnt't have anything related to this
	prometheus.Register(memoryCadvisorPgfaults)     // done for both pgmajfault and pgfault

	// Register cpu metrics with prometheus
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuNrPeriods)
	prometheus.MustRegister(cpuNrThrottled)
	prometheus.MustRegister(cpuThrottled)

	// Register IO metrics with prometheus
	prometheus.MustRegister(ioRbytes)
	prometheus.MustRegister(ioWbytes)
	prometheus.MustRegister(ioRios)
	prometheus.MustRegister(ioWios)

	// Register IO metrics with prometheus
	prometheus.MustRegister(ioCadvisorRbytes)
	prometheus.MustRegister(ioCadvisorWbytes)
	prometheus.MustRegister(ioCadvisorRios)
	prometheus.MustRegister(ioCadvisorWios)
}

func cgroupMetrics(hasMemoryController, hasCPUController, hasIOController bool, cadvisorMetrics bool) {
	blockDevices()

	for {
		cgItems := cgServices()

		for _, item := range cgItems {
			if hasMemoryController {
				go cgroupMemoryMetrics(item, cadvisorMetrics)
			}

			if hasCPUController {
				go cgroupCPUMetrics(item)
			}

			if hasIOController {
				go cgroupIOMetrics(item, cadvisorMetrics)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
