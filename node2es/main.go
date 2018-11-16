package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hongxincn/promexp/node2es/add"
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
		configFile       = kingpin.Flag("config.file", "Configuration file").Default("node2es_config.yml").String()
		keyToBeEncrypted = kingpin.Flag("encrypting.key", "encrypting inputted key, and exist immediately").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(fmt.Sprintf("%s\n%s\n%s", VERSION, BUILD_TIME, GO_VERSION))
	kingpin.Parse()

	if *keyToBeEncrypted != "" {
		fmt.Printf("Encrypted String\n%s\n", config.EncryptingString(*keyToBeEncrypted))
		os.Exit(0)
	}

	err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = es.NewEsClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := time.Tick(time.Duration(60) * time.Second)
	go func() {
		for {
			prometheus.LoadMetrics()
			<-c
		}
	}()
	fmt.Println("监听端口", config.Config.Listen_on)
	if config.Config.Add_fields.Api_url != "" {
		addFieldsEndPoint := add.NewAddFields()
		http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
			addFieldsEndPoint.SetReloadFlag()
			w.Write([]byte(`{"acknowledged":"true"}`))
		})
		http.HandleFunc("/-/instancesMapping", func(w http.ResponseWriter, r *http.Request) {
			w.Write(addFieldsEndPoint.GetInstancesMapping())
		})
	}
	err = http.ListenAndServe(config.Config.Listen_on, nil)
	if err != nil {
		fmt.Println(err)
	}
}
