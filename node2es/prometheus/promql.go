package prometheus

var (
	prom_queries = map[string]string{
		"cpu_usage": `100 - (avg by (instance) (irate(node_cpu_seconds_total{ mode="idle"}[5m])) * 100)`,
		"mem_usage": `((node_memory_MemTotal_bytes - (node_memory_MemAvailable_bytes or (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)))/node_memory_MemTotal_bytes)*100`,
	}
)

func getPromQueries() map[string]string {
	return prom_queries
}
