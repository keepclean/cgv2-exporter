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

func cgroupsMetrics() {
	hasMemoryController := hasController("memory")
	hasCPUController := hasController("cpu")

	go func() {
		for {
			cgItems := cgServices()

			for _, item := range cgItems {
				if hasMemoryController {
					stat := &memoryStat{}
					if err := parseMemoryStat(item, stat); err != nil {
						log.Println(err)
					}
					parseMemoryFiles(item, stat)

					memoryAnon.WithLabelValues(item).Set(float64(stat.Anon))
					memoryFile.WithLabelValues(item).Set(float64(stat.File))
					memoryKernelStack.WithLabelValues(item).Set(float64(stat.KernelStack))
					memorySlab.WithLabelValues(item).Set(float64(stat.Slab))
					memorySock.WithLabelValues(item).Set(float64(stat.Sock))
					memoryShmem.WithLabelValues(item).Set(float64(stat.Shmem))
					memoryFileMapped.WithLabelValues(item).Set(float64(stat.FileMapped))
					memoryFileDirty.WithLabelValues(item).Set(float64(stat.FileDirty))
					memoryFileWriteback.WithLabelValues(item).Set(float64(stat.FileWriteback))
					memoryInactiveAnon.WithLabelValues(item).Set(float64(stat.InactiveAnon))
					memoryActiveAnon.WithLabelValues(item).Set(float64(stat.ActiveAnon))
					memoryInactiveFile.WithLabelValues(item).Set(float64(stat.InactiveFile))
					memoryActiveFile.WithLabelValues(item).Set(float64(stat.ActiveFile))
					memoryUnevictable.WithLabelValues(item).Set(float64(stat.Unevictable))
					memorySlabReclaimable.WithLabelValues(item).Set(float64(stat.SlabReclaimable))
					memorySlabUnreclaimable.WithLabelValues(item).Set(float64(stat.SlabUnreclaimable))
					memoryPgfault.WithLabelValues(item).Set(float64(stat.Pgfault))
					memoryPgmajfault.WithLabelValues(item).Set(float64(stat.Pgmajfault))
					memoryPgrefill.WithLabelValues(item).Set(float64(stat.Pgrefill))
					memoryPgscan.WithLabelValues(item).Set(float64(stat.Pgscan))
					memoryPgsteal.WithLabelValues(item).Set(float64(stat.Pgsteal))
					memoryPgactivate.WithLabelValues(item).Set(float64(stat.Pgactivate))
					memoryPgdeactivate.WithLabelValues(item).Set(float64(stat.Pgdeactivate))
					memoryPglazyfree.WithLabelValues(item).Set(float64(stat.Pglazyfree))
					memoryPglazyfreed.WithLabelValues(item).Set(float64(stat.Pglazyfreed))
					memoryWorkingsetRefault.WithLabelValues(item).Set(float64(stat.WorkingsetRefault))
					memoryWorkingsetActivate.WithLabelValues(item).Set(float64(stat.WorkingsetActivate))
					memoryWorkingsetNodereclaim.WithLabelValues(item).Set(float64(stat.WorkingsetNodereclaim))
					memoryCurrent.WithLabelValues(item).Set(float64(stat.Current))
					memoryHigh.WithLabelValues(item).Set(float64(stat.High))
					memoryLow.WithLabelValues(item).Set(float64(stat.Low))
					memoryMax.WithLabelValues(item).Set(float64(stat.Max))
					memoryMin.WithLabelValues(item).Set(float64(stat.Min))
				}

				if hasCPUController {
					stat := &cpuStat{}
					if err := parseCPUStat(item, stat); err != nil {
						log.Println(err)
					}

					cpuUsage.WithLabelValues(item).Set(stat.Usage)
					cpuUser.WithLabelValues(item).Set(stat.User)
					cpuSystem.WithLabelValues(item).Set(stat.System)
					cpuNrPeriods.WithLabelValues(item).Set(stat.NrPeriods)
					cpuNrThrottled.WithLabelValues(item).Set(stat.NrThrottled)
					cpuThrottled.WithLabelValues(item).Set(stat.Throttled)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}
