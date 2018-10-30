package prometheus

var (
	prom_queries = map[string]string{
		"cpu_usage":     `100 - (avg by (instance) (irate(node_cpu_seconds_total{ mode="idle"}[5m])) * 100)`,
		"mem_usage":     `((node_memory_MemTotal_bytes - (node_memory_MemAvailable_bytes or (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)))/node_memory_MemTotal_bytes)*100`,
		"swap_usage":    `(node_memory_SwapCached_bytes/node_memory_SwapTotal_bytes) * 100`,
		"procs_blocked": `node_procs_blocked`,
		"procs_running": `node_procs_running`,
		"load1":         `node_load1`,
	}
)

func getPromQueries() map[string]string {
	return prom_queries
}
