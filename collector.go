package main

import "github.com/prometheus/client_golang/prometheus"

var (
	memoryAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_anon_count",
			Help: "Amount of memory used in anonymous mappings such as brk(), sbrk(), and mmap(MAP_ANONYMOUS)",
		},
		[]string{"service"},
	)
	memoryFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_count",
			Help: "Amount of memory used to cache filesystem data including tmpfs and shared memory.",
		},
		[]string{"service"},
	)
	memoryKernelStack = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_kernel_stack_count",
			Help: "Amount of memory allocated to kernel stacks.",
		},
		[]string{"service"},
	)
	memorySlab = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_count",
			Help: "Amount of memory used for storing in-kernel data structures.",
		},
		[]string{"service"},
	)
	memorySock = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_sock_count",
			Help: "Amount of memory used in network transmission buffers",
		},
		[]string{"service"},
	)
	memoryShmem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_shmem_count",
			Help: "Amount of cached filesystem data that is swap-backed, such as tmpfs, shm segments, shared anonymous mmap()s",
		},
		[]string{"service"},
	)
	memoryFileMapped = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_mapped_count",
			Help: "Amount of cached filesystem data mapped with mmap()",
		},
		[]string{"service"},
	)
	memoryFileDirty = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_dirty_count",
			Help: "Amount of cached filesystem data that was modified but not yet written back to disk",
		},
		[]string{"service"},
	)
	memoryFileWriteback = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_file_writeback_count",
			Help: "Amount of cached filesystem data that was modified and is currently being written back to disk",
		},
		[]string{"service"},
	)
	memoryInactiveAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_inactive_anon_count",
			Help: "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryActiveAnon = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_active_anon_count",
			Help: "Amount of swap-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryInactiveFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_inactive_file_count",
			Help: "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryActiveFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_active_file_count",
			Help: "Amount of filesystem-backed memory on the internal memory management lists used by the page reclaim algorithm",
		},
		[]string{"service"},
	)
	memoryUnevictable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_unevictable_count",
			Help: "About of memory which never will be reclaimed from memory",
		},
		[]string{"service"},
	)
	memorySlabReclaimable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_reclaimable_count",
			Help: "Part of slab that might be reclaimed, such as dentries and inodes.",
		},
		[]string{"service"},
	)
	memorySlabUnreclaimable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_slab_unreclaimable_count",
			Help: "Part of slab that cannot be reclaimed on memory pressure.",
		},
		[]string{"service"},
	)
	memoryPgfault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgfault_count",
			Help: "Total number of page faults incurred",
		},
		[]string{"service"},
	)
	memoryPgmajfault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgmajfault_count",
			Help: "Number of major page faults incurred",
		},
		[]string{"service"},
	)
	memoryPgrefill = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgrefill_count",
			Help: "Amount of scanned pages (in an active LRU list)",
		},
		[]string{"service"},
	)
	memoryPgscan = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgscan_count",
			Help: "Amount of scanned pages (in an inactive LRU list)",
		},
		[]string{"service"},
	)
	memoryPgsteal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgsteal_count",
			Help: "Amount of reclaimed pages",
		},
		[]string{"service"},
	)
	memoryPgactivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgactivate_count",
			Help: "Amount of pages moved to the active LRU list",
		},
		[]string{"service"},
	)
	memoryPgdeactivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pgdeactivate_count",
			Help: "Amount of pages moved to the inactive LRU list",
		},
		[]string{"service"},
	)
	memoryPglazyfree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pglazyfree_count",
			Help: "Amount of pages postponed to be freed under memory pressure",
		},
		[]string{"service"},
	)
	memoryPglazyfreed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_pglazyfreed_count",
			Help: "Amount of reclaimed lazyfree pages",
		},
		[]string{"service"},
	)
	memoryWorkingsetRefault = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_refault_count",
			Help: "Number of refaults of previously evicted pages",
		},
		[]string{"service"},
	)
	memoryWorkingsetActivate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_activate_count",
			Help: "Number of refaulted pages that were immediately activated",
		},
		[]string{"service"},
	)
	memoryWorkingsetNodereclaim = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_workingset_nodereclaim_count",
			Help: "Number of times a shadow node has been reclaimed",
		},
		[]string{"service"},
	)
)

func init() {
	// Register metrics with prometheus
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
}

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
}
