package main

import (
	"log"
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

	// Register cpu metrics with prometheus
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuNrPeriods)
	prometheus.MustRegister(cpuNrThrottled)
	prometheus.MustRegister(cpuThrottled)
}

func cgroupsMetrics() {
	hasMemory := hasController("memory")
	hasCPU := hasController("cpu")

	go func() {
		for {
			cgItems := cgServices()
			memStats := make(map[string]memoryStat)
			cpuStats := make(map[string]cpuStat)

			for _, item := range cgItems {
				if hasMemory {
					stat := &memoryStat{}
					if err := parseMemoryStat(item, stat); err != nil {
						log.Fatalln(err)
					}
					memStats[item] = *stat
				}

				if hasCPU {
					stat := &cpuStat{}
					if err := parseCPUStat(item, stat); err != nil {
						log.Fatalln(err)
					}
					cpuStats[item] = *stat
				}
			}

			for _, item := range cgItems {
				if hasMemory {
					memoryAnon.WithLabelValues(item).Set(float64(memStats[item].Anon))
					memoryFile.WithLabelValues(item).Set(float64(memStats[item].File))
					memoryKernelStack.WithLabelValues(item).Set(float64(memStats[item].KernelStack))
					memorySlab.WithLabelValues(item).Set(float64(memStats[item].Slab))
					memorySock.WithLabelValues(item).Set(float64(memStats[item].Sock))
					memoryShmem.WithLabelValues(item).Set(float64(memStats[item].Shmem))
					memoryFileMapped.WithLabelValues(item).Set(float64(memStats[item].FileMapped))
					memoryFileDirty.WithLabelValues(item).Set(float64(memStats[item].FileDirty))
					memoryFileWriteback.WithLabelValues(item).Set(float64(memStats[item].FileWriteback))
					memoryInactiveAnon.WithLabelValues(item).Set(float64(memStats[item].InactiveAnon))
					memoryActiveAnon.WithLabelValues(item).Set(float64(memStats[item].ActiveAnon))
					memoryInactiveFile.WithLabelValues(item).Set(float64(memStats[item].InactiveFile))
					memoryActiveFile.WithLabelValues(item).Set(float64(memStats[item].ActiveFile))
					memoryUnevictable.WithLabelValues(item).Set(float64(memStats[item].Unevictable))
					memorySlabReclaimable.WithLabelValues(item).Set(float64(memStats[item].SlabReclaimable))
					memorySlabUnreclaimable.WithLabelValues(item).Set(float64(memStats[item].SlabUnreclaimable))
					memoryPgfault.WithLabelValues(item).Set(float64(memStats[item].Pgfault))
					memoryPgmajfault.WithLabelValues(item).Set(float64(memStats[item].Pgmajfault))
					memoryPgrefill.WithLabelValues(item).Set(float64(memStats[item].Pgrefill))
					memoryPgscan.WithLabelValues(item).Set(float64(memStats[item].Pgscan))
					memoryPgsteal.WithLabelValues(item).Set(float64(memStats[item].Pgsteal))
					memoryPgactivate.WithLabelValues(item).Set(float64(memStats[item].Pgactivate))
					memoryPgdeactivate.WithLabelValues(item).Set(float64(memStats[item].Pgdeactivate))
					memoryPglazyfree.WithLabelValues(item).Set(float64(memStats[item].Pglazyfree))
					memoryPglazyfreed.WithLabelValues(item).Set(float64(memStats[item].Pglazyfreed))
					memoryWorkingsetRefault.WithLabelValues(item).Set(float64(memStats[item].WorkingsetRefault))
					memoryWorkingsetActivate.WithLabelValues(item).Set(float64(memStats[item].WorkingsetActivate))
					memoryWorkingsetNodereclaim.WithLabelValues(item).Set(float64(memStats[item].WorkingsetNodereclaim))
				}

				if hasCPU {
					cpuUsage.WithLabelValues(item).Set(cpuStats[item].Usage)
					cpuUser.WithLabelValues(item).Set(cpuStats[item].User)
					cpuSystem.WithLabelValues(item).Set(cpuStats[item].System)
					cpuNrPeriods.WithLabelValues(item).Set(cpuStats[item].NrPeriods)
					cpuNrThrottled.WithLabelValues(item).Set(cpuStats[item].NrThrottled)
					cpuThrottled.WithLabelValues(item).Set(cpuStats[item].Throttled)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}
