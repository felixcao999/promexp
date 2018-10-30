package main

import (
	"fmt"
	"os"

	"github.com/hongxincn/promexp/node2es/config"
	"github.com/hongxincn/promexp/node2es/es"
	"github.com/hongxincn/promexp/node2es/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	var (
		configFile = kingpin.Flag("config.file", "Configuration file").Default("node2es_config.yml").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(fmt.Sprintf("%s\n%s\n%s", VERSION, BUILD_TIME, GO_VERSION))
	kingpin.Parse()

	err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = es.NewEs5Client()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prometheus.LoadMetrics()
}
