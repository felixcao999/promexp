package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prometheus/common/model"

	"github.com/hongxincn/promexp/node2es/config"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
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
	queries := getPromQueries()
	nms := NewNodeMetrics()
	for k, v := range queries {
		r, err := query_api.Query(ctx, v, time.Now())
		if err != nil {
			fmt.Printf("Error occuried while querying, err: %v, query: %s \n", err, v)
		}
		v, ok := r.(model.Vector)
		if ok {
			for _, vr := range v {
				nms.Set(fmt.Sprintf("%v", vr.Metric["instance"]), k, vr.Value.String())
			}
		}
	}
	for k, v := range nms.metrics {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			fmt.Println("json marshal error, key=%s, value=%v", k, v)
		}
		fmt.Println(string(jsonBytes))
	}
}
