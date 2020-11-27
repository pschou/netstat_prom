# netstat_prom

A simple netstat like exporter for prometheus to track ongoing connections.

To build:
```
make
```

To run:
```
netstat_prom
```

Program flags:
```
# ./netstat_prom -h
Usage of ./netstat_prom:
  -listen string
        ip and port to listen on (default ":9733")
  -showbins
        show time bin counts (default true)
  -showtcp
        show all tcp connections
  -showudp
        show all udp connections
  -tcpbins string
        timebins in seconds for tcp connections (default "120,300,600,900,1800,3600,5400,7200,9000,10800,12600,14400,28800,57600,86400")
```
