package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hongxincn/promexp/node2es/config"
	"github.com/hongxincn/promexp/node2es/es"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func LoadMetrics() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url := config.Config.Prometheus.Url
	client, err := api.NewClient(api.Config{Address: url})
	if err != nil {
		fmt.Printf("Error when trying to connect to %s, err: %v \n", url, err)
	}

	query_api := v1.NewAPI(client)
	nms := NewNodeMetrics()
	ip_port_label := config.Config.Promql.Ip_port_label
	var wg = sync.WaitGroup{}
	getMetric := func(ctx context.Context, q config.PromQuery, qt time.Time) {
		defer wg.Done()
		r, err := query_api.Query(ctx, q.Query, qt)
		if err != nil {
			fmt.Printf("Error occuried while querying, err: %v, query: %s \n", err, q.Query)
		}
		v, ok := r.(model.Vector)
		if ok {
			for _, vr := range v {
				if q.Keep_labels {
					nms.Add(string(vr.Metric[model.LabelName(ip_port_label)]), q.Metric, vr.Metric, (float64)(vr.Value), ip_port_label)
				} else {
					nms.Set(string(vr.Metric[model.LabelName(ip_port_label)]), q.Metric, (float64)(vr.Value))
				}
			}
		}
	}
	qt := time.Now()
	for _, q := range config.Config.Promql.Querys {
		wg.Add(1)
		go getMetric(ctx, q, qt)
	}
	wg.Wait()
	t1 := time.Now().Unix()
	index := es.GetIndex(t1)
	count := 0
	for k, v := range nms.metrics {
		vi := map[string]interface{}{}
		for k2, v2 := range v {
			vi[k2] = v2
		}
		mlv := nms.metrics_with_labels[k]
		for k3, v3 := range mlv {
			vi[k3] = v3
		}
		vi["instance_id"] = k
		vi["instance_ip"] = k[:strings.LastIndex(k, ":")]
		vi["instance_port"] = k[strings.LastIndex(k, ":")+1:]
		vi["timestamp"] = t1

		now := time.Now()
		local1, err := time.LoadLocation("") //same as "UTC"
		if err != nil {
			fmt.Println(err)
		}
		sTimeProcessed := now.In(local1).Format("2006-01-02T15:04:05.000Z")
		vi["@timestamp"] = sTimeProcessed

		jsonBytes, err := json.Marshal(vi)
		if err != nil {
			fmt.Printf("json marshal error, key=%s, value=%v \n", k, vi)
		} else {
			es.Client.AddBulkRequest(index, string(jsonBytes))
			count++
		}
		if count >= 2000 {
			es.Client.SubmitBulkRequest()
			count = 0
		}
	}
	es.Client.SubmitBulkRequest()
}
