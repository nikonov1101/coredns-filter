package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coredns/caddy"

	// make all plugins available
	_ "github.com/coredns/coredns/plugin/any"
	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/cancel"
	_ "github.com/coredns/coredns/plugin/chaos"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/dnssec"
	_ "github.com/coredns/coredns/plugin/dnstap"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/hosts"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/metrics"
	_ "github.com/coredns/coredns/plugin/minimal"
	_ "github.com/coredns/coredns/plugin/nsid"
	_ "github.com/coredns/coredns/plugin/pprof"
	_ "github.com/coredns/coredns/plugin/tls"
	_ "github.com/coredns/coredns/plugin/trace"

	"gitlab.com/nikonov1101/coredns-filter/blocklist"
	"gitlab.com/nikonov1101/coredns-filter/piholedb"
)

var (
	flagDNSForwarder1        = flag.String("dns.forward1", "1.1.1.1", "Forwarding DNS #1")
	flagDNSForwarder2        = flag.String("dns.forward2", "1.0.0.1", "Forwarding DNS #2")
	flagPrometheusListenAddr = flag.String("prom.listen", "0.0.0.0:9999", "Prometheus listen address")
	flagGravityDBPath        = flag.String("db", "gravity.db", "path to PiHole's gravity.db")
	flagHostsFilePath        = flag.String("hosts-file", "", "path to a /etc/hosts-like file, disabled if empty")
)

func main() {
	flag.Parse()

	rd, err := piholedb.New(*flagGravityDBPath)
	if err != nil {
		log.Printf("failed to read blocklist db: %v", err)
		os.Exit(1)
	}

	if err := blocklist.InitPlugin(rd); err != nil {
		log.Printf("failed to init blocklist plugin: %v", err)
		os.Exit(2)
	}

	caddyFile := generateCaddyFile()
	instance, err := caddy.Start(caddyFile)
	if err != nil {
		log.Printf("failed to start DNS server: %v", err)
		os.Exit(3)
	}

	log.Println("caddy: waiting...")
	instance.Wait()
}

// generateCaddyFile turns given flags into a Caddyfile for CoreDNS.
func generateCaddyFile() caddy.CaddyfileInput {
	// TODO(nikonov): how to TLS forwarders?
	forwardServers := []string{
		*flagDNSForwarder1,
		*flagDNSForwarder2,
	}

	opts := []string{
		"prometheus " + *flagPrometheusListenAddr,
		"log",
		"errors",
		"blocklist",
		"forward . " + strings.Join(forwardServers, " "),
		"cache",
	}

	if flagHostsFilePath != nil && len(*flagHostsFilePath) > 0 {
		opts = append(opts, "hosts "+*flagHostsFilePath)
	}

	bs := ". {\n" + strings.Join(opts, "\n") + "\n}"
	log.Println(">>>>>> starting with config:")
	fmt.Println(bs)
	log.Println("<<<<<< end config")

	return caddy.CaddyfileInput{
		Filepath:       "",
		Contents:       []byte(bs),
		ServerTypeName: "dns",
	}
}
