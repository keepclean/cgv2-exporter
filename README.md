# cgv2-exporter

cgv2-exporter - is a simple Prometheus exporter for unified cgroup

Currently supports cpu, io, memory controllers.

```
$ ./cgv2-exporter --help
Usage of ./cgv2-exporter:
  -cadvisor_metrics
        Add to exported metrics cadvisor style metrics
  -listen_ip string
        IP to listen on, defaults to all IPs
  -port int
        port to listen (default 8888)
  -scraping_interval uint
        Scraping interval in seconds (default 5)
```
