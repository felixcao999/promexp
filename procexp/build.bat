set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-w -s -X main.VERSION=v0.0.1 -X 'main.BUILD_TIME=2019-02-20 15:30' -X 'main.GO_VERSION=1.11'"