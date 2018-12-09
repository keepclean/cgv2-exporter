package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const cgDir string = "/sys/fs/cgroup/system.slice/"

func main() {
	var argIP = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
	var argPort = flag.Int("port", 8888, "port to listen")
	flag.Parse()

	cgItems := cgServices()
	stats := make(map[string]memoryStat)

	for _, item := range cgItems {
		stat := &memoryStat{}

		if err := memStat(item, stat); err != nil {
			log.Fatalln(err)
		}

		stats[item] = *stat
	}

	for _, item := range cgItems {
		memoryAnon.WithLabelValues(item).Set(float64(stats[item].Anon))
		memoryFile.WithLabelValues(item).Set(float64(stats[item].File))
		memoryKernelStack.WithLabelValues(item).Set(float64(stats[item].KernelStack))
		memorySlab.WithLabelValues(item).Set(float64(stats[item].Slab))
		memorySock.WithLabelValues(item).Set(float64(stats[item].Sock))
		memoryShmem.WithLabelValues(item).Set(float64(stats[item].Shmem))
		memoryFileMapped.WithLabelValues(item).Set(float64(stats[item].FileMapped))
		memoryFileDirty.WithLabelValues(item).Set(float64(stats[item].FileDirty))
		memoryFileWriteback.WithLabelValues(item).Set(float64(stats[item].FileWriteback))
		memoryInactiveAnon.WithLabelValues(item).Set(float64(stats[item].InactiveAnon))
		memoryActiveAnon.WithLabelValues(item).Set(float64(stats[item].ActiveAnon))
		memoryInactiveFile.WithLabelValues(item).Set(float64(stats[item].InactiveFile))
		memoryActiveFile.WithLabelValues(item).Set(float64(stats[item].ActiveFile))
		memoryUnevictable.WithLabelValues(item).Set(float64(stats[item].Unevictable))
		memorySlabReclaimable.WithLabelValues(item).Set(float64(stats[item].SlabReclaimable))
		memorySlabUnreclaimable.WithLabelValues(item).Set(float64(stats[item].SlabUnreclaimable))
		memoryPgfault.WithLabelValues(item).Set(float64(stats[item].Pgfault))
		memoryPgmajfault.WithLabelValues(item).Set(float64(stats[item].Pgmajfault))
		memoryPgrefill.WithLabelValues(item).Set(float64(stats[item].Pgrefill))
		memoryPgscan.WithLabelValues(item).Set(float64(stats[item].Pgscan))
		memoryPgsteal.WithLabelValues(item).Set(float64(stats[item].Pgsteal))
		memoryPgactivate.WithLabelValues(item).Set(float64(stats[item].Pgactivate))
		memoryPgdeactivate.WithLabelValues(item).Set(float64(stats[item].Pgdeactivate))
		memoryPglazyfree.WithLabelValues(item).Set(float64(stats[item].Pglazyfree))
		memoryPglazyfreed.WithLabelValues(item).Set(float64(stats[item].Pglazyfreed))
		memoryWorkingsetRefault.WithLabelValues(item).Set(float64(stats[item].WorkingsetRefault))
		memoryWorkingsetActivate.WithLabelValues(item).Set(float64(stats[item].WorkingsetActivate))
		memoryWorkingsetNodereclaim.WithLabelValues(item).Set(float64(stats[item].WorkingsetNodereclaim))
	}

	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("%s:%d", *argIP, *argPort)
	log.Println("Starting web server on: ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func cgServices() (items []string) {
	entries, err := ioutil.ReadDir(cgDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range entries {
		if !item.IsDir() {
			continue
		}
		if !strings.HasSuffix(item.Name(), ".service") {
			continue
		}
		items = append(items, item.Name())
	}

	return
}

func memStat(item string, stat *memoryStat) error {
	file, err := os.Open(filepath.Join(cgDir, item, "memory.stat"))
	if err != nil {
		return err
	}
	defer close(file)

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

	stat.Anon += raw["anon"]
	stat.File += raw["file"]
	stat.KernelStack += raw["kernel_stack"]
	stat.Slab += raw["slab"]
	stat.Sock += raw["sock"]
	stat.Shmem += raw["shmem"]
	stat.FileMapped += raw["file_mapped"]
	stat.FileDirty += raw["file_dirty"]
	stat.FileWriteback += raw["file_writeback"]
	stat.InactiveAnon += raw["inactive_anon"]
	stat.ActiveAnon += raw["active_anon"]
	stat.InactiveFile += raw["inactive_file"]
	stat.ActiveFile += raw["active_file"]
	stat.Unevictable += raw["unevictable"]
	stat.SlabReclaimable += raw["slab_reclaimable"]
	stat.SlabUnreclaimable += raw["slab_unreclaimable"]
	stat.Pgfault += raw["pgfault"]
	stat.Pgmajfault += raw["pgmajfault"]
	stat.Pgrefill += raw["pgrefill"]
	stat.Pgscan += raw["pgscan"]
	stat.Pgsteal += raw["pgsteal"]
	stat.Pgactivate += raw["pgactivate"]
	stat.Pgdeactivate += raw["pgdeactivate"]
	stat.Pglazyfree += raw["pglazyfree"]
	stat.Pglazyfreed += raw["pglazyfreed"]
	stat.WorkingsetRefault += raw["workingset_refault"]
	stat.WorkingsetActivate += raw["workingset_activate"]
	stat.WorkingsetNodereclaim += raw["workingset_nodereclaim"]

	return nil
}

func parseKV(s string) (string, uint64, error) {
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return "", 0, errors.New("Invalid format")
	}

	v, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return fields[0], v, nil

}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
