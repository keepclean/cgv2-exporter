package main

import (
	"time"
)

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
