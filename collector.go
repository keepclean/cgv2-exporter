package main

import (
	"log"
	"time"
)

func cgroupMetrics(cadvisorMetrics bool, scrapingInterval uint) {
	blockDevices()
	controllers := cgroupControllers()

	for {
		services, err := systemdServices()
		if err != nil {
			log.Println("An error has happned while discovering systemd services", err)
		}

		for _, service := range services {
			if controllers["memory"] {
				go cgroupMemoryMetrics(service, cadvisorMetrics)
			}

			if controllers["cpu"] {
				go cgroupCPUMetrics(service)
			}

			if controllers["io"] {
				go cgroupIOMetrics(service, cadvisorMetrics)
			}
		}

		time.Sleep(time.Duration(scrapingInterval) * time.Second)
	}
}
