package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	argIP           string
	argPort         uint
	cadvisorMetrics bool
	cgDir           = "/sys/fs/cgroup/system.slice/"
)

func init() {
	flag.StringVar(&argIP, "listen_ip", "", "IP to listen on, defaults to all IPs")
	flag.UintVar(&argPort, "port", 8888, "port to listen")
	flag.BoolVar(&cadvisorMetrics, "cadvisor_metrics", false, "Add to exported metrics cadvisor in style metrics")
	flag.Parse()
}

func main() {
	var srv http.Server
	idleConnsClosed := make(chan struct{})
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
		<-osSignal

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("HTTP server shutdown:", err)
		}
		close(idleConnsClosed)
	}()

	blockDevices()

	exporter, err := newExporter(cpuMetrics, ioMetrics, memoryMetrics)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
             <head><title>Exporter for unified cgroup of systemd services</title></head>
             <body>
             <h1>Exporter for unified cgroup of systemd services</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
        </html>`))
		if err != nil {
			log.Println(err)
		}
	})

	srv.Addr = fmt.Sprintf("%s:%d", argIP, argPort)
	log.Println("Starting web server on: ", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("ListenAndServe:", err)
	}

	<-idleConnsClosed
}
