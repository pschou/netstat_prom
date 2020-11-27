
build:
	GOPATH=${CURDIR} CGO_ENABLED=0 go build -o netstat_prom main.go
