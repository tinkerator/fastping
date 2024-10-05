# fastping

fastping is a Go language ICMP ping library, inspired by the `AnyEvent::FastPing`
Perl module, for quickly sending ICMP ECHO REQUEST packets. Original Perl module
is available at http://search.cpan.org/~mlehmann/AnyEvent-FastPing-2.01/

All original functions haven't been implemented yet.

[![GoDoc](https://pkg.go.dev/zappem.net/pub/net/fastping?status.svg)](https://pkg.go.dev/zappem.net/pub/net/fastping)

## Installation

Install and update with `go get zappem.net/pub/net/fastping`.

## Examples

Import this package and write

```go
p := fastping.NewPinger()
ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
p.AddIPAddr(ra)
p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
	fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
}
p.OnIdle = func() {
	fmt.Println("finish")
}
err = p.Run()
if err != nil {
	fmt.Println(err)
}
```

The example sends an ICMP packet and waits for a response. If it receives a
response, it calls the "receive" callback. After that, once MaxRTT time has
passed, it calls the "idle" callback. For more details,
refer [to the godoc][godoc], and if you need more examples,
please see "cmd/ping/ping.go".

**Caution** This package implements ICMP ping using both raw socket
  and UDP. If your program needs to use this package in raw socket
  mode, it needs to be privileged. Either run it as a `root` user or,
  under Linux, with the `CAP_NET_RAW` capability enabled. For that
  latter case, try the following:

```
$ go build cmd/ping/ping.go
$ sudo setcap cap_net_raw=ep ./ping
$ ./ping localhost
```

## License
fastping is under MIT License. See the [LICENSE][license] file for details.

[godoc]: https://pkg.go.dev/zappem.net/pub/net/fastping
[license]: https://github.com/tinkerator/fastping/blob/main/LICENSE

## History

This package was originally (2024) cloned from
https://github.com/tatsushid/go-fastping.