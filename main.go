package main

import (
	"flag"
	"log"
	"os"

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

	// init our custom plugin
	_ "gitlab.com/nikonov1101/coredns-filter/blocklist"
)

func main() {
	flag.Parse()

	caddyFile, err := caddy.LoadCaddyfile("dns")
	if err != nil {
		log.Printf("failed to load caddyfile: %v", err)
		os.Exit(1)
	}

	instance, err := caddy.Start(caddyFile)
	if err != nil {
		log.Printf("failed to start DNS server: %v", err)
		os.Exit(2)
	}

	log.Println("caddy: running...")
	instance.Wait()
}
