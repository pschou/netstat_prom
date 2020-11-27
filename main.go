package main

import (
	"./github.com/drael/GOnetstat"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var time_bins []int
var show_bins bool
var show_tcp bool
var show_udp bool

func main() {

	listen := flag.String("listen", ":9733", "ip and port to listen on")
	bins := flag.String("tcpbins", "120,300,600,900,1800,3600,5400,7200,9000,10800,12600,14400,28800,57600,86400", "timebins in seconds for tcp connections")
	showtcp := flag.Bool("showtcp", false, "show all tcp connections")
	showudp := flag.Bool("showudp", false, "show all udp connections")
	showbins := flag.Bool("showbins", true, "show time bin counts")
	flag.Parse()
	show_tcp = *showtcp
	show_udp = *showudp
	show_bins = *showbins
	for _, i := range strings.Split(*bins, ",") {
		val, err := strconv.Atoi(i)
		if err == nil {
			time_bins = append(time_bins, val)
		}
	}
	time_bins = append(time_bins, 1e10)

	http.HandleFunc("/metrics", Collect)
	http.ListenAndServe(*listen, nil)
}
func Collect(w http.ResponseWriter, req *http.Request) {
	// Get Udp data, you can use GOnetstat.Tcp() to get TCP data
	d := GOnetstat.Tcp6()
	now := time.Now()
	var time_counts = make([]int, len(time_bins))

	for _, p := range d {
		t := -float64(p.Mtime.Sub(now)) / 1e9
		if show_tcp {
			fmt.Fprintf(w, "netstat_info{proto=\"tcp6\",localIP=\"%v\",localPort=\"%v\",remoteIP=\"%v\",remotePort=\"%v\",pid=\"%v\",user=\"%v\",state=\"%v\",exe=\"%v\"} %v\n",
				p.Ip, p.Port, p.ForeignIp, p.ForeignPort, p.Pid, p.User, p.State, p.Exe, t)
		}
		if show_bins {
			if p.State != "LISTEN" {
				for k, v := range time_bins {
					if int(t) < v {
						time_counts[k]++
						break
					}
				}
			}
		}
	}

	d = GOnetstat.Tcp()
	now = time.Now()

	for _, p := range d {
		t := -float64(p.Mtime.Sub(now)) / 1e9
		if show_tcp {
			fmt.Fprintf(w, "netstat_info{proto=\"tcp\",localIP=\"%v\",localPort=\"%v\",remoteIP=\"%v\",remotePort=\"%v\",pid=\"%v\",user=\"%v\",state=\"%v\",exe=\"%v\"} %v\n",
				p.Ip, p.Port, p.ForeignIp, p.ForeignPort, p.Pid, p.User, p.State, p.Exe, t)
		}
		if show_bins {
			if p.State != "LISTEN" {
				for k, v := range time_bins {
					if int(t) < v {
						time_counts[k]++
						break
					}
				}
			}
		}
	}

	for k, v := range time_bins {
		if v == 1e10 {
			fmt.Fprintf(w, "netstat_time_bins{proto=\"tcp\",le=\"inf\"} %v\n", time_counts[k])
		} else {
			fmt.Fprintf(w, "netstat_time_bins{proto=\"tcp\",le=\"%v\"} %v\n", v, time_counts[k])
		}
	}

	if show_udp {
		d = GOnetstat.Udp6()
		now = time.Now()
		for _, p := range d {
			t := -float64(p.Mtime.Sub(now)) / 1e9
			fmt.Fprintf(w, "netstat_info{proto=\"udp6\",localIP=\"%v\",localPort=\"%v\",remoteIP=\"%v\",remotePort=\"%v\",pid=\"%v\",user=\"%v\",state=\"%v\",exe=\"%v\"} %v\n",
				p.Ip, p.Port, p.ForeignIp, p.ForeignPort, p.Pid, p.User, p.State, p.Exe, t)
		}
		d = GOnetstat.Udp()
		now = time.Now()
		for _, p := range d {
			t := -float64(p.Mtime.Sub(now)) / 1e9
			fmt.Fprintf(w, "netstat_info{proto=\"udp\",localIP=\"%v\",localPort=\"%v\",remoteIP=\"%v\",remotePort=\"%v\",pid=\"%v\",user=\"%v\",state=\"%v\",exe=\"%v\"} %v\n",
				p.Ip, p.Port, p.ForeignIp, p.ForeignPort, p.Pid, p.User, p.State, p.Exe, t)
		}
	}
}
