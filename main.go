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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const cgDir string = "/sys/fs/cgroup/system.slice/"

var (
	argIP            string
	argPort          int
	cadvisorMetrics  bool
	scrapingInterval uint
)

func init() {
	flag.StringVar(&argIP, "listen_ip", "", "IP to listen on, defaults to all IPs")
	flag.IntVar(&argPort, "port", 8888, "port to listen")
	flag.BoolVar(&cadvisorMetrics, "cadvisor_metrics", false, "Add to exported metrics cadvisor style metrics")
	flag.UintVar(&scrapingInterval, "scraping_interval", 5, "Scraping interval in seconds")
	flag.Parse()
}

func main() {
	var srv http.Server
	idleConnsClosed := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("HTTP server shutdown:", err)
		}
		close(idleConnsClosed)
	}()

	go cgroupMetrics(
		hasController("memory"),
		hasController("cpu"),
		hasController("io"),
		cadvisorMetrics,
		scrapingInterval,
	)

	http.Handle("/metrics", promhttp.Handler())

	srv.Addr = fmt.Sprintf("%s:%d", argIP, argPort)
	log.Println("Starting web server on: ", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("ListenAndServe:", err)
	}

	<-idleConnsClosed
}
