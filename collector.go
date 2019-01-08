package main

import (
	"time"
)

func cgroupMetrics(hasMemoryController, hasCPUController, hasIOController, cadvisorMetrics bool, scrapingInterval uint) {
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

		time.Sleep(time.Duration(scrapingInterval) * time.Second)
	}
}
