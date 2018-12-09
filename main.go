package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const cgDir string = "/sys/fs/cgroup/system.slice/"

func main() {
	var argIP = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
	var argPort = flag.Int("port", 8888, "port to listen")
	flag.Parse()

	cgroupsMetrics()

	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("%s:%d", *argIP, *argPort)
	log.Println("Starting web server on: ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
