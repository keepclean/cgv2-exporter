package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type cpuStat struct {
	Usage       float64
	User        float64
	System      float64
	NrPeriods   float64
	NrThrottled float64
	Throttled   float64
}

var (
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_usage_seconds_total",
			Help: "Cumulative cpu time consumed",
		},
		[]string{"service"},
	)
	cpuUser = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_user_seconds_total",
			Help: "Cumulative user cpu time consumed",
		},
		[]string{"service"},
	)
	cpuSystem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_system_seconds_total",
			Help: "Cumulative system cpu time consumed",
		},
		[]string{"service"},
	)
	cpuNrPeriods = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_nr_periods_total",
			Help: "Number of enforcement intervals that have elapsed.",
		},
		[]string{"service"},
	)
	cpuNrThrottled = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_nr_throttled_periods_total",
			Help: "Number of times the group has been throttled/limited.",
		},
		[]string{"service"},
	)
	cpuThrottled = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_cpu_throttled_seconds_total",
			Help: "The total time duration for which entities of the group have been throttled.",
		},
		[]string{"service"},
	)
)

func cgroupCPUMetrics(service string) {
	stat := &cpuStat{}
	if err := parseCPUStat(service, stat); err != nil {
		log.Println(err)
	}

	cpuUsage.WithLabelValues(service).Set(stat.Usage)
	cpuUser.WithLabelValues(service).Set(stat.User)
	cpuSystem.WithLabelValues(service).Set(stat.System)
	cpuNrPeriods.WithLabelValues(service).Set(stat.NrPeriods)
	cpuNrThrottled.WithLabelValues(service).Set(stat.NrThrottled)
	cpuThrottled.WithLabelValues(service).Set(stat.Throttled)
}

func parseCPUStat(service string, stat *cpuStat) error {
	file, err := os.Open(filepath.Join(cgDir, service, "cpu.stat"))
	if err != nil {
		return err
	}
	defer closeFile(file)

	raw := make(map[string]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, err := parseKV(scanner.Text())
		if err != nil {
			return err
		}

		if strings.HasSuffix(key, "_usec") {
			raw[key] = float64(value) / 1e9
			continue
		}

		raw[key] = float64(value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	stat.Usage = raw["usage_usec"]
	stat.User = raw["user_usec"]
	stat.System = raw["system_usec"]
	stat.NrPeriods = raw["nr_periods"]
	stat.NrThrottled = raw["nr_throttled"]
	stat.Throttled = raw["throttled_usec"]

	return nil
}
