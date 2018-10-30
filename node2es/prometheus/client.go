package prometheus

import (
	"context"
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

	client, err := api.NewClient(api.Config{Address: config.Config.Prometheus.Url})
	if err != nil {
		fmt.Println("Ooops connecting to the API", err)
	}

	query_api := v1.NewAPI(client)

	queries := getPromQueries()

	for k, v := range queries {
		r, _ := query_api.Query(ctx, v, time.Now())

		v, ok := r.(model.Vector)
		if ok {
			fmt.Printf("%v\n", k)
			for _, vr := range v {
				fmt.Printf("%v\n", vr.Metric["instance"])
				fmt.Printf("%v\n", vr.Value)
			}
		}
	}
}
