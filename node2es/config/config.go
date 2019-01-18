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
	Listen_on  string
	Prometheus struct {
		Url string
	}
	Es struct {
		Urls         []string
		Username     string
		Password     string
		Index        string
		Version      float32
		Default_type string
	}
	Promql struct {
		Instance_id struct {
			Label, Regex, Replacement string
			Is_ip_port                bool
		}
		Ip_port_label string
		Querys        []PromQuery
	}
	Add_fields struct {
		Api_url string
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
	if Config.Promql.Instance_id.Regex == "" {
		Config.Promql.Instance_id.Regex = "(.*)"
	}
	if Config.Promql.Instance_id.Replacement == "" {
		Config.Promql.Instance_id.Replacement = "$1"
	}
	if Config.Es.Default_type == "" {
		if Config.Es.Version < 6 {
			Config.Es.Default_type = "doc"
		} else {
			Config.Es.Default_type = "_doc"
		}
	}

	return nil
}
