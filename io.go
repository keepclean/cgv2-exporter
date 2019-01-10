package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type ioStat struct {
	Rbytes float64
	Wbytes float64
	Rios   float64
	Wios   float64
}

var (
	// io.stat file
	ioRbytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_io_read_bytes",
			Help: "Bytes read",
		},
		[]string{"service", "device"},
	)
	ioWbytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_io_write_bytes",
			Help: "Bytes written",
		},
		[]string{"service", "device"},
	)
	ioRios = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_io_read_operations",
			Help: "Number of read IOs",
		},
		[]string{"service", "device"},
	)
	ioWios = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_io_write_operations",
			Help: "Number of write IOs",
		},
		[]string{"service", "device"},
	)

	// cadvisor style memory metrics for the backward compability
	ioCadvisorRbytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_fs_reads_bytes_total",
			Help: "Cumulative count of bytes read",
		},
		[]string{"app_name", "device"},
	)
	ioCadvisorWbytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_fs_writes_bytes_total",
			Help: "Cumulative count of bytes written",
		},
		[]string{"app_name", "device"},
	)
	ioCadvisorRios = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_fs_reads_total",
			Help: "Cumulative count of reads completed",
		},
		[]string{"app_name", "device"},
	)
	ioCadvisorWios = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_fs_writes_total",
			Help: "Cumulative count of writes completed",
		},
		[]string{"app_name", "device"},
	)
)

func cgroupIOMetrics(service string, cadvisorMetrics bool) {
	stat, err := parseIOStat(service)
	if err != nil {
		log.Println(err)
	}

	for d, s := range stat {
		ioRbytes.WithLabelValues(service, d).Set(s["rbytes"])
		ioWbytes.WithLabelValues(service, d).Set(s["wbytes"])
		ioRios.WithLabelValues(service, d).Set(s["rios"])
		ioWios.WithLabelValues(service, d).Set(s["wios"])

		if cadvisorMetrics {
			d = fmt.Sprint("/dev/", d)
			ioCadvisorRbytes.WithLabelValues(service, d).Set(s["rbytes"])
			ioCadvisorWbytes.WithLabelValues(service, d).Set(s["wbytes"])
			ioCadvisorRios.WithLabelValues(service, d).Set(s["rios"])
			ioCadvisorWios.WithLabelValues(service, d).Set(s["wios"])
		}
	}
}

func parseIOStat(service string) (map[string]map[string]float64, error) {
	file, err := os.Open(filepath.Join(cgDir, service, "io.stat"))
	if err != nil {
		return map[string]map[string]float64{}, err
	}
	defer closeFile(file)

	raw := make(map[string]map[string]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 5 {
			return map[string]map[string]float64{}, errors.New("Invalid io.stat file format")
		}

		device := devices[fields[0]]
		for _, substring := range fields[1:] {
			if raw[device] == nil {
				raw[device] = map[string]float64{}
			}
			kv := strings.Split(substring, "=")
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				v = float64(0)
			}
			raw[device][kv[0]] = v
		}
	}

	if err := scanner.Err(); err != nil {
		return map[string]map[string]float64{}, err
	}

	return raw, nil
}
