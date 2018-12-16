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

func main() {
	var argIP = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
	var argPort = flag.Int("port", 8888, "port to listen")
	flag.Parse()

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

	cgroupsMetrics()

	http.Handle("/metrics", promhttp.Handler())

	srv.Addr = fmt.Sprintf("%s:%d", *argIP, *argPort)
	log.Println("Starting web server on: ", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("ListenAndServe:", err)
	}

	<-idleConnsClosed
}
