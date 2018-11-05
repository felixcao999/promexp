package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var Config Node2EsConfig

const (
	ENCRYPT_KEY = "node2es_encryption_key"
)

type PromQuery struct {
	Metric      string
	Query       string
	Keep_labels bool
}

type Node2EsConfig struct {
	Prometheus struct {
		Url string
	}
	Es struct {
		Urls     []string
		Username string
		Password string
		Index    string
	}
	Promql struct {
		Ip_port_label string
		Querys        []PromQuery
	}
}

// LoadConfig loads the specified YAML configuration file.
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		return err
	}

	return nil
}
