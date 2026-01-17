package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"zappem.net/pub/net/fastping"
)

var (
	debug = flag.Bool("debug", false, "use to enable fastping package debug logging")
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

func main() {
	var useUDP bool
	flag.BoolVar(&useUDP, "udp", false, "use non-privileged datagram-oriented UDP as ICMP endpoints")
	flag.BoolVar(&useUDP, "u", false, "use non-privileged datagram-oriented UDP as ICMP endpoints (shorthand)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] hostname [source]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	hostname := flag.Arg(0)
	if len(hostname) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	source := ""
	if flag.NArg() > 1 {
		source = flag.Arg(1)
	}

	p := fastping.NewPinger()
	if useUDP {
		p.Network("udp")
	}
	p.Debug = *debug

	netProto := "ip4:icmp"
	if strings.Index(hostname, ":") != -1 {
		netProto = "ip6:ipv6-icmp"
	}
	ra, err := net.ResolveIPAddr(netProto, hostname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if source != "" {
		p.Source(source)
	}

	results := make(map[string]*response)
	results[ra.String()] = nil
	p.AddIPAddr(ra)

	onRecv, onIdle := make(chan *response), make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.MaxRTT = time.Second
	p.RunLoop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	var good, bad int64
loop:
	for {
		select {
		case <-c:
			fmt.Println("interrupted")
			break loop
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			for host, r := range results {
				if r == nil {
					fmt.Printf("%s : unreachable %v\n", host, time.Now())
					bad++
				} else {
					fmt.Printf("%s : %v %v\n", host, r.rtt, time.Now())
					good++
				}
				results[host] = nil
			}
		case <-p.Done():
			if err = p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
			}
			break loop
		}
	}
	signal.Stop(c)
	p.Stop()

	den := float64(good+bad) + .0001
	log.Printf("attempts=%.0f, success=%d(%.1f%%), failure=%d(%.1f%%)", den, good, 100*float64(good)/den, bad, 100*float64(bad)/den)
}
