package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type memoryStat struct {
	Anon                  uint64
	File                  uint64
	KernelStack           uint64
	Slab                  uint64
	Sock                  uint64
	Shmem                 uint64
	FileMapped            uint64
	FileDirty             uint64
	FileWriteback         uint64
	InactiveAnon          uint64
	ActiveAnon            uint64
	InactiveFile          uint64
	ActiveFile            uint64
	Unevictable           uint64
	SlabReclaimable       uint64
	SlabUnreclaimable     uint64
	Pgfault               uint64
	Pgmajfault            uint64
	Pgrefill              uint64
	Pgscan                uint64
	Pgsteal               uint64
	Pgactivate            uint64
	Pgdeactivate          uint64
	Pglazyfree            uint64
	Pglazyfreed           uint64
	WorkingsetRefault     uint64
	WorkingsetActivate    uint64
	WorkingsetNodereclaim uint64
	Current               uint64
	High                  uint64
	Low                   uint64
	Max                   uint64
	Min                   uint64
}

var (
	memoryAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_anon_bytes",
			Help: "Amount of memory used in anonymous mappings such as brk(), sbrk(), and mmap(MAP_ANONYMOUS)",
		},
		[]string{"service"},
	)
	memoryFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_bytes",
			Help: "Amount of memory used to cache filesystem data including tmpfs and shared memory.",
		},
		[]string{"service"},
	)
	memoryKernelStack = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_kernel_stack_bytes",
			Help: "Amount of memory allocated to kernel stacks.",
		},
		[]string{"service"},
	)
	memorySlab = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_bytes",
			Help: "Amount of memory used for storing in-kernel data structures.",
		},
		[]string{"service"},
	)
	memorySock = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_sock_bytes",
			Help: "Amount of memory used in network transmission buffers",
		},
		[]string{"service"},
	)
	memoryShmem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_shmem_bytes",
			Help: "Amount of cached filesystem data that is swap-backed, such as tmpfs, shm segments, shared anonymous mmap()s",
		},
		[]string{"service"},
	)
	memoryFileMapped = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_mapped_bytes",
			Help: "Amount of cached filesystem data mapped with mmap()",
		},
		[]string{"service"},
	)
	memoryFileDirty = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_dirty_bytes",
			Help: "Amount of cached filesystem data that was modified but not yet written back to disk",
		},
		[]string{"service"},
	)
	memoryFileWriteback = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_writeback_bytes",
			Help: "Amount of cached filesystem data that was modified and is currently being written back to disk",
		},
		[]string{"service"},
	)
	memoryInactiveAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_inactive_anon_bytes",
			Help: "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryActiveAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_active_anon_bytes",
			Help: "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryInactiveFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_inactive_file_bytes",
			Help: "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryActiveFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_active_file_bytes",
			Help: "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryUnevictable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_unevictable_bytes",
			Help: "About of memory which never will be reclaimed from memory",
		},
		[]string{"service"},
	)
	memorySlabReclaimable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_reclaimable_bytes",
			Help: "Part of slab that might be reclaimed, such as dentries and inodes.",
		},
		[]string{"service"},
	)
	memorySlabUnreclaimable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_unreclaimable_bytes",
			Help: "Part of slab that cannot be reclaimed on memory pressure.",
		},
		[]string{"service"},
	)
	memoryPgfault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgfault_pages",
			Help: "Total number of page faults incurred",
		},
		[]string{"service"},
	)
	memoryPgmajfault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgmajfault_pages",
			Help: "Number of major page faults incurred",
		},
		[]string{"service"},
	)
	memoryPgrefill = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgrefill_pages",
			Help: "Amount of scanned pages (in an active LRU list)",
		},
		[]string{"service"},
	)
	memoryPgscan = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgscan_pages",
			Help: "Amount of scanned pages (in an inactive LRU list)",
		},
		[]string{"service"},
	)
	memoryPgsteal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgsteal_pages",
			Help: "Amount of reclaimed pages",
		},
		[]string{"service"},
	)
	memoryPgactivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgactivate_pages",
			Help: "Amount of pages moved to the active LRU list",
		},
		[]string{"service"},
	)
	memoryPgdeactivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgdeactivate_pages",
			Help: "Amount of pages moved to the inactive LRU list",
		},
		[]string{"service"},
	)
	memoryPglazyfree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pglazyfree_pages",
			Help: "Amount of pages postponed to be freed under memory pressure",
		},
		[]string{"service"},
	)
	memoryPglazyfreed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pglazyfreed_pages",
			Help: "Amount of reclaimed lazyfree pages",
		},
		[]string{"service"},
	)
	memoryWorkingsetRefault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_refault_pages",
			Help: "Number of refaults of previously evicted pages",
		},
		[]string{"service"},
	)
	memoryWorkingsetActivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_activate_pages",
			Help: "Number of refaulted pages that were immediately activated",
		},
		[]string{"service"},
	)
	memoryWorkingsetNodereclaim = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_nodereclaim_pages",
			Help: "Number of times a shadow node has been reclaimed",
		},
		[]string{"service"},
	)
	memoryCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_current_bytes",
			Help: "The total amount of memory currently being used by the cgroup and its descendants.",
		},
		[]string{"service"},
	)
	memoryHigh = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_high_bytes",
			Help: "Memory usage throttle limit.",
		},
		[]string{"service"},
	)
	memoryLow = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_low_bytes",
			Help: "Best-effort memory protection.",
		},
		[]string{"service"},
	)
	memoryMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_max_bytes",
			Help: "Memory usage hard limit.",
		},
		[]string{"service"},
	)
	memoryMin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_min_bytes",
			Help: "Hard memory protection.",
		},
		[]string{"service"},
	)
)

func cgroupMemoryMetics(item string) {
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

func parseMemoryStat(item string, stat *memoryStat) error {
	file, err := os.Open(filepath.Join(cgDir, item, "memory.stat"))
	if err != nil {
		return err
	}
	defer closeFile(file)

	raw := make(map[string]uint64)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, err := parseKV(scanner.Text())
		if err != nil {
			return err
		}
		raw[key] = value
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	stat.Anon = raw["anon"]
	stat.File = raw["file"]
	stat.KernelStack = raw["kernel_stack"]
	stat.Slab = raw["slab"]
	stat.Sock = raw["sock"]
	stat.Shmem = raw["shmem"]
	stat.FileMapped = raw["file_mapped"]
	stat.FileDirty = raw["file_dirty"]
	stat.FileWriteback = raw["file_writeback"]
	stat.InactiveAnon = raw["inactive_anon"]
	stat.ActiveAnon = raw["active_anon"]
	stat.InactiveFile = raw["inactive_file"]
	stat.ActiveFile = raw["active_file"]
	stat.Unevictable = raw["unevictable"]
	stat.SlabReclaimable = raw["slab_reclaimable"]
	stat.SlabUnreclaimable = raw["slab_unreclaimable"]
	stat.Pgfault = raw["pgfault"]
	stat.Pgmajfault = raw["pgmajfault"]
	stat.Pgrefill = raw["pgrefill"]
	stat.Pgscan = raw["pgscan"]
	stat.Pgsteal = raw["pgsteal"]
	stat.Pgactivate = raw["pgactivate"]
	stat.Pgdeactivate = raw["pgdeactivate"]
	stat.Pglazyfree = raw["pglazyfree"]
	stat.Pglazyfreed = raw["pglazyfreed"]
	stat.WorkingsetRefault = raw["workingset_refault"]
	stat.WorkingsetActivate = raw["workingset_activate"]
	stat.WorkingsetNodereclaim = raw["workingset_nodereclaim"]

	return nil
}

var totalRAM = totalRAMMemory()

func parseMemoryFiles(item string, stat *memoryStat) {
	memoryFiles, err := controllerFiles("memory", item)
	if err != nil {
		log.Println(err)
		return
	}

	raw := make(map[string]uint64)
	for _, f := range memoryFiles {
		// memory.stat file is parsed in parseMemoryStat func
		if f == "memory.stat" || f == "memory.events" {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(cgDir, item, f))
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.Contains(string(file), "max") {
			raw[f] = totalRAM
			continue
		}

		v, err := strconv.ParseUint(strings.TrimSuffix(string(file), "\n"), 10, 64)
		if err != nil {
			v = 0
		}
		raw[f] = v
	}

	stat.Current = raw["memory.current"]
	stat.High = raw["memory.high"]
	stat.Low = raw["memory.low"]
	stat.Max = raw["memory.max"]
	stat.Min = raw["memory.min"]
}
