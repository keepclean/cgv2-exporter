package main

import (
	"bufio"
	"errors"
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
	"github.com/sirupsen/logrus"
)

const cgDir string = "/sys/fs/cgroup/system.slice/"

/*
var fooMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "foo_metric", Help: "Shows whether a foo has occurred in out cluster"})

var barMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "bar_metric", Help: "Shows whether a bar has occurred in out cluster"})

func init() {
	// Register metrics with prometheus
	prometheus.MustRegister(fooMetric)
	prometheus.MustRegister(barMetric)

	// Set fooMetric to 1
	fooMetric.Set(1)
	// Ste barMetric to 0
	barMetric.Set(0)
}
*/

func main() {

	cgItems, err := ioutil.ReadDir(cgDir)
	if err != nil {
		log.Fatal(err)
	}
	stats := make(map[string]memoryStat)

	for _, item := range cgItems {
		if item.IsDir() && strings.HasSuffix(item.Name(), ".service") {

			stat := &memoryStat{}

			if err := memStat(item.Name(), stat); err != nil {
				log.Fatalln(err)
			}

			stats[item.Name()] = *stat
		}
	}

	for k, v := range stats {
		fmt.Printf("%s: %+v\n\n", k, v)
	}

	http.Handle("/metrics", promhttp.Handler())
	logrus.Info("Beginning to sever on port :8000")
	logrus.Fatal(http.ListenAndServe(":8000", nil))

}

type memoryStat struct {
	anon                   uint64
	file                   uint64
	kernel_stack           uint64
	slab                   uint64
	sock                   uint64
	shmem                  uint64
	file_mapped            uint64
	file_dirty             uint64
	file_writeback         uint64
	inactive_anon          uint64
	active_anon            uint64
	inactive_file          uint64
	active_file            uint64
	unevictable            uint64
	slab_reclaimable       uint64
	slab_unreclaimable     uint64
	pgfault                uint64
	pgmajfault             uint64
	pgrefill               uint64
	pgscan                 uint64
	pgsteal                uint64
	pgactivate             uint64
	pgdeactivate           uint64
	pglazyfree             uint64
	pglazyfreed            uint64
	workingset_refault     uint64
	workingset_activate    uint64
	workingset_nodereclaim uint64
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

	stat.anon += raw["anon"]
	stat.file += raw["file"]
	stat.kernel_stack += raw["kernel_stack"]
	stat.slab += raw["slab"]
	stat.sock += raw["sock"]
	stat.shmem += raw["shmem"]
	stat.file_mapped += raw["file_mapped"]
	stat.file_dirty += raw["file_dirty"]
	stat.file_writeback += raw["file_writeback"]
	stat.inactive_anon += raw["inactive_anon"]
	stat.active_anon += raw["active_anon"]
	stat.inactive_file += raw["inactive_file"]
	stat.active_file += raw["active_file"]
	stat.unevictable += raw["unevictable"]
	stat.slab_reclaimable += raw["slab_reclaimable"]
	stat.slab_unreclaimable += raw["slab_unreclaimable"]
	stat.pgfault += raw["pgfault"]
	stat.pgmajfault += raw["pgmajfault"]
	stat.pgrefill += raw["pgrefill"]
	stat.pgscan += raw["pgscan"]
	stat.pgsteal += raw["pgsteal"]
	stat.pgactivate += raw["pgactivate"]
	stat.pgdeactivate += raw["pgdeactivate"]
	stat.pglazyfree += raw["pglazyfree"]
	stat.pglazyfreed += raw["pglazyfreed"]
	stat.workingset_refault += raw["workingset_refault"]
	stat.workingset_activate += raw["workingset_activate"]
	stat.workingset_nodereclaim += raw["workingset_nodereclaim"]

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
