package main

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "container"
)

var (
	memoryLabelNames = []string{"app_name"}
	cpuLabelNames    = []string{"app_name"}
	ioLabelNames     = []string{"app_name", "device"}

	cpuMetrics = map[string]*prometheus.Desc{
		"usage_usec":     newMetric("cpu", "usage_seconds_total", "Cumulative cpu time consumed", cpuLabelNames),
		"user_usec":      newMetric("cpu", "user_seconds_total", "Cumulative user cpu time consumed", cpuLabelNames),
		"system_usec":    newMetric("cpu", "system_seconds_total", "Cumulative system cpu time consumed", cpuLabelNames),
		"nr_periods":     newMetric("cpu", "nr_periods_total", "Number of enforcement intervals that have elapsed.", cpuLabelNames),
		"nr_throttled":   newMetric("cpu", "nr_throttled_periods_total", "Number of times the group has been throttled/limited.", cpuLabelNames),
		"throttled_usec": newMetric("cpu", "throttled_seconds_total", "The total time duration for which entities of the group have been throttled.", cpuLabelNames),
	}
	ioMetrics = map[string]*prometheus.Desc{
		"rbytes": newMetric("io", "read_bytes", "Bytes read", ioLabelNames),
		"wbytes": newMetric("io", "write_bytes", "Bytes written", ioLabelNames),
		"rios":   newMetric("io", "read_operations", "Number of read IOs", ioLabelNames),
		"wios":   newMetric("io", "write_operations", "Number of write IOs", ioLabelNames),
		"dbytes": newMetric("io", "discarded_bytes", "Number of bytes discarded", ioLabelNames),
		"dios":   newMetric("io", "discarded_operations", "Number of discard or trim IOs.", ioLabelNames),
	}
	cadvisorIOMetrics = map[string]*prometheus.Desc{
		"rbytes": newMetric("fs", "reads_bytes_total", "Cumulative count of bytes read", ioLabelNames),
		"wbytes": newMetric("fs", "writes_bytes_total", "Cumulative count of bytes written", ioLabelNames),
		"rios":   newMetric("fs", "reads_total", "Cumulative count of reads completed", ioLabelNames),
		"wios":   newMetric("fs", "writes_total", "Cumulative count of writes completed", ioLabelNames),
		"dbytes": newMetric("fs", "discarded_bytes_total", "Cumulative count of bytes discarded", ioLabelNames),
		"dios":   newMetric("fs", "discarded_total", "Cumulative count of discard or trim IOs.", ioLabelNames),
	}
	memoryMetrics = map[string]*prometheus.Desc{
		"anon":                   newMetric("memory", "anon_bytes", "Amount of memory used in anonymous mappings such as brk(), sbrk(), and mmap(MAP_ANONYMOUS)", memoryLabelNames),
		"file":                   newMetric("memory", "file_bytes", "Amount of memory used to cache filesystem data including tmpfs and shared memory.", memoryLabelNames),
		"kernel_stack":           newMetric("memory", "kernel_stack_bytes", "Amount of memory allocated to kernel stacks.", memoryLabelNames),
		"slab":                   newMetric("memory", "slab_bytes", "Amount of memory used for storing in-kernel data structures.", memoryLabelNames),
		"sock":                   newMetric("memory", "sock_bytes", "Amount of memory used in network transmission buffers", memoryLabelNames),
		"shmem":                  newMetric("memory", "shmem_bytes", "Amount of cached filesystem data that is swap-backed, such as tmpfs, shm segments, shared anonymous mmap()s", memoryLabelNames),
		"file_mapped":            newMetric("memory", "file_mapped_bytes", "Amount of cached filesystem data mapped with mmap()", memoryLabelNames),
		"file_dirty":             newMetric("memory", "file_dirty_bytes", "Amount of cached filesystem data that was modified but not yet written back to disk", memoryLabelNames),
		"file_writeback":         newMetric("memory", "file_writeback_bytes", "Amount of cached filesystem data that was modified and is currently being written back to disk", memoryLabelNames),
		"inactive_anon":          newMetric("memory", "inactive_anon_bytes", "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm", memoryLabelNames),
		"active_anon":            newMetric("memory", "active_anon_bytes", "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm", memoryLabelNames),
		"inactive_file":          newMetric("memory", "inactive_file_bytes", "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm", memoryLabelNames),
		"active_file":            newMetric("memory", "active_file_bytes", "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm", memoryLabelNames),
		"unevictable":            newMetric("memory", "unevictable_bytes", "About of memory which never will be reclaimed from memory", memoryLabelNames),
		"slab_reclaimable":       newMetric("memory", "slab_reclaimable_bytes", "Part of slab that might be reclaimed, such as dentries and inodes.", memoryLabelNames),
		"slab_unreclaimable":     newMetric("memory", "slab_unreclaimable_bytes", "Part of slab that cannot be reclaimed on memory pressure.", memoryLabelNames),
		"pgfault":                newMetric("memory", "pgfault_pages", "Total number of page faults incurred", memoryLabelNames),
		"pgmajfault":             newMetric("memory", "pgmajfault_pages", "Number of major page faults incurred", memoryLabelNames),
		"pgrefill":               newMetric("memory", "pgrefill_pages", "Amount of scanned pages (in an active LRU list)", memoryLabelNames),
		"pgscan":                 newMetric("memory", "pgscan_pages", "Amount of scanned pages (in an inactive LRU list)", memoryLabelNames),
		"pgsteal":                newMetric("memory", "pgsteal_pages", "Amount of reclaimed pages", memoryLabelNames),
		"pgactivate":             newMetric("memory", "pgactivate_pages", "Amount of pages moved to the active LRU list", memoryLabelNames),
		"pgdeactivate":           newMetric("memory", "pgdeactivate_pages", "Amount of pages moved to the inactive LRU list", memoryLabelNames),
		"pglazyfree":             newMetric("memory", "pglazyfree_pages", "Amount of pages postponed to be freed under memory pressure", memoryLabelNames),
		"pglazyfreed":            newMetric("memory", "pglazyfreed_pages", "Amount of reclaimed lazyfree pages", memoryLabelNames),
		"workingset_refault":     newMetric("memory", "workingset_refault_pages", "Number of refaults of previously evicted pages", memoryLabelNames),
		"workingset_activate":    newMetric("memory", "workingset_activate_pages", "Number of refaulted pages that were immediately activated", memoryLabelNames),
		"workingset_nodereclaim": newMetric("memory", "workingset_nodereclaim_pages", "Number of times a shadow node has been reclaimed", memoryLabelNames),
		"memory.current":         newMetric("memory", "current_bytes", "The total amount of memory currently being used by the cgroup and its descendants.", memoryLabelNames),
		"memory.high":            newMetric("memory", "high_bytes", "Memory usage throttle limit.", memoryLabelNames),
		"memory.low":             newMetric("memory", "low_bytes", "Best-effort memory protection.", memoryLabelNames),
		"memory.max":             newMetric("memory", "max_bytes", "Memory usage hard limit.", memoryLabelNames),
		"memory.min":             newMetric("memory", "min_bytes", "Hard memory protection.", memoryLabelNames),
		"low":                    newMetric("memory", "low_events", "The number of times the cgroup is reclaimed due to high memory pressure even though its usage is under the low boundary.", memoryLabelNames),
		"high":                   newMetric("memory", "high_events", "The number of times processes of the cgroup are throttled and routed to perform direct memory reclaim because the high memory boundary was exceeded.", memoryLabelNames),
		"max":                    newMetric("memory", "max_events", "The number of times the cgroup's memory usage was about to go over the max boundary.", memoryLabelNames),
		"oom":                    newMetric("memory", "oom_events", "The number of time the cgroup's memory usage was reached the limit and allocation was about to fail.", memoryLabelNames),
		"oom_kill":               newMetric("memory", "oom_kill_events", "The number of processes belonging to this cgroup killed by any kind of OOM killer.", memoryLabelNames),
	}
	cadvisorMemMetrics = map[string]*prometheus.Desc{
		"cache":                       newMetric("memory", "cache", "Number of bytes of page cache memory.", memoryLabelNames),
		"failcnt":                     newMetric("memory", "failcnt", "Number of memory usage hits limits.", memoryLabelNames),
		"max_usage":                   newMetric("memory", "max_usage_bytes", "Maximum memory usage recorded in bytes.", memoryLabelNames),
		"usage":                       newMetric("memory", "usage_bytes", "Current memory usage in bytes, including all memory regardless of when it was accessed.", memoryLabelNames),
		"rss":                         newMetric("memory", "rss", "Size of RSS in bytes.", memoryLabelNames),
		"swap":                        newMetric("memory", "swap", "Container swap usage in bytes.", memoryLabelNames),
		"working_set":                 newMetric("memory", "working_set_bytes", "Current working set in bytes.", memoryLabelNames),
		"container_spec_memory_limit": newMetric("memory", "container_spec_memory_limit_bytes", "Memory limit for the container.", memoryLabelNames),
		"container_spec_memory_reservation_limit": newMetric("memory", "container_spec_memory_reservation_limit_bytes", "Memory reservation limit for the container.", memoryLabelNames),
		"container_spec_memory_swap_limit":        newMetric("memory", "container_spec_memory_swap_limit_bytes", "Memory swap limit for the container.", memoryLabelNames),
		"failures":                                newMetric("memory", "failures_total", "Cumulative count of memory allocation failures.", []string{"app_name", "scope", "type"}),
	}
)

func newMetric(controller, metricName, docString string, controllerLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, controller, metricName), docString, controllerLabels, nil)
}

// Exporter collects unified cgroups stats for systemd sevices
// and exports them using the prometheus metrics package.
type Exporter struct {
	mutex                                 sync.RWMutex
	cpuMetrics, ioMetrics, memoryMetrics  map[string]*prometheus.Desc
	cadvisorIOMetrics, cadvisorMemMetrics map[string]*prometheus.Desc
}

// Describe describes all the metrics ever exported by the cgv2-exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range cpuMetrics {
		ch <- m
	}
	for _, m := range ioMetrics {
		ch <- m
	}
	for _, m := range memoryMetrics {
		ch <- m
	}

	if cadvisorMetrics {
		for _, m := range cadvisorIOMetrics {
			ch <- m
		}
		for _, m := range cadvisorMemMetrics {
			ch <- m
		}
	}
}

// Collect fetches the stats systemd services and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	services, err := systemdServices()
	if err != nil {
		log.Println("An error has happned while discovering systemd services", err)
	}

	for _, service := range services {
		cgroupMetrics(service, ch)
	}
}

// newExporter returns an initialized Exporter.
func newExporter(cpuMetrics, ioMetrics, cgroupMetrics map[string]*prometheus.Desc) (*Exporter, error) {
	return &Exporter{
		cpuMetrics:         cpuMetrics,
		ioMetrics:          ioMetrics,
		cadvisorIOMetrics:  cadvisorIOMetrics,
		memoryMetrics:      memoryMetrics,
		cadvisorMemMetrics: cadvisorMemMetrics,
	}, nil
}

func cgroupMetrics(service string, ch chan<- prometheus.Metric) {
	files, err := cgroupFiles(service)
	if err != nil {
		log.Println(err)
	}

	serviceStats := make(map[string]float64)
	serviceIOStats := make(map[string]map[string]float64)

	for _, f := range files {
		switch f {
		case "cpu.stat":
			if err := parseCPUKvFile(service, f, serviceStats); err != nil {
				log.Println(err)
			}
		case "memory.stat", "memory.events":
			if err := parseMemoryKvFile(service, f, serviceStats); err != nil {
				log.Println(err)
			}
		case "memory.current", "memory.high", "memory.low", "memory.max", "memory.min":
			if err := parseMemoryFile(service, f, serviceStats); err != nil {
				log.Println(err)
			}
		case "io.stat":
			if err := parseIOKvFile(service, f, serviceIOStats); err != nil {
				log.Println(err)
			}
		}
	}

	for name, metric := range cpuMetrics {
		ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, serviceStats[name], service)
	}
	for device, stats := range serviceIOStats {
		for name, metric := range ioMetrics {
			ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, stats[name], service, device)
			if cadvisorMetrics {
				ch <- prometheus.MustNewConstMetric(cadvisorIOMetrics[name], prometheus.GaugeValue, stats[name], service, device)
			}
		}
	}
	for name, metric := range memoryMetrics {
		ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, serviceStats[name], service)
	}
	if cadvisorMetrics {
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["cache"], prometheus.GaugeValue, serviceStats["file"], service)
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["failcnt"], prometheus.GaugeValue, serviceStats["max_events"], service)
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["usage"], prometheus.GaugeValue, serviceStats["memory.current"], service)
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["rss"], prometheus.GaugeValue, serviceStats["anon"], service)

		var workingSet float64
		if !(serviceStats["memory.current"] < serviceStats["inactive_file"]) {
			workingSet = serviceStats["memory.current"] - serviceStats["inactive_file"]
		}
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["working_set"], prometheus.GaugeValue, workingSet, service)

		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["container_spec_memory_limit"], prometheus.GaugeValue, serviceStats["memory.max"], service)
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["failures"], prometheus.GaugeValue, serviceStats["pgfault"], service, "container", "pgfault")
		ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["failures"], prometheus.GaugeValue, serviceStats["pgmajfault"], service, "container", "pgmajfault")

		// TODO some metrics don't have relative unified cgroup "analogue"
		// ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["swap"], prometheus.GaugeValue, serviceStats["XXX"], service)
		// ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["max_usage"], prometheus.GaugeValue, serviceStats["XXX"], service)
		// ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["container_spec_memory_reservation_limit"], prometheus.GaugeValue, serviceStats["XXX"], service)
		// ch <- prometheus.MustNewConstMetric(cadvisorMemMetrics["container_spec_memory_swap_limit"], prometheus.GaugeValue, serviceStats["XXX"], service)
	}
}
