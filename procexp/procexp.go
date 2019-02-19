package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/shirou/gopsutil/process"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "proc"
	clientID  = "process_exporter"
)

var (
	VERSION           string
	BUILD_TIME        string
	GO_VERSION        string
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scrape_success"),
		"process_exporter: scrape succeed or not.",
		nil,
		nil,
	)
)

type Exporter struct {
}

func init() {
	prometheus.MustRegister(version.NewCollector(clientID))
}

func NewExporter() (*Exporter, error) {
	return &Exporter{}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeSuccessDesc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ps, err := process.Processes()
	if err != nil {
		fmt.Println(err)
		ch <- prometheus.MustNewConstMetric(
			scrapeSuccessDesc, prometheus.GaugeValue, float64(0))
		return
	}
	for _, proc := range ps {
		cmdLine, err := proc.Cmdline()
		if err != nil {
			fmt.Println(err)
		}
		procName, _ := proc.Name()
		if cmdLine == "" {
			cmdLine = procName
		}
		//proc.Pid
		userName, err := proc.Username()
		if err != nil {
			fmt.Println(err)
		}

		cpuPercent, _ := proc.CPUPercent()
		memPercent, _ := proc.MemoryPercent()
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "cpu_percent"),
				"Process CPU Percent",
				[]string{"pid", "user", "cmd"}, nil,
			),
			prometheus.GaugeValue, float64(cpuPercent), strconv.Itoa(int(proc.Pid)), userName, cmdLine)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "mem_percent"),
				"Process Memory Percent",
				[]string{"pid", "user", "cmd"}, nil,
			),
			prometheus.GaugeValue, float64(memPercent), strconv.Itoa(int(proc.Pid)), userName, cmdLine)

	}
	ch <- prometheus.MustNewConstMetric(
		scrapeSuccessDesc, prometheus.GaugeValue, float64(1))
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9256").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(fmt.Sprintf("process exporter version %s by hx\nbuild on %s\nGO %s\n", VERSION, BUILD_TIME, GO_VERSION))

	kingpin.Parse()

	log.Println("Starting process_exporter, listening on ", *listenAddress)

	exporter, err := NewExporter()
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	        <head><title>Process Exporter</title></head>
	        <body>
	        <h1>Process Exporter</h1>
	        <p><a href='` + *metricsPath + `'>Metrics</a></p>
	        </body>
	        </html>`))
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
