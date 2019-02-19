module github.com/hongxincn/promexp/procexp

replace golang.org/x/sync v0.0.0-20181108010431-42b317875d0f => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f

replace golang.org/x/crypto v0.0.0-20180904163835-0709b304e793 => github.com/golang/crypto v0.0.0-20180904163835-0709b304e793

replace (
	golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33 => github.com/golang/sys v0.0.0-20180905080454-ebe1bf3edb33
	golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5 => github.com/golang/sys v0.0.0-20181116152217-5ac8a444bdc5
)

replace (
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
	golang.org/x/net v0.0.0-20181201002055-351d144fa1fc => github.com/golang/net v0.0.0-20181201002055-351d144fa1fc
)

require (
	github.com/StackExchange/wmi v0.0.0-20181212234831-e0a55b97c705 // indirect
	github.com/go-ole/go-ole v1.2.2 // indirect
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/common v0.2.0
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
