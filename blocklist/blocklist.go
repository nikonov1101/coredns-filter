package blocklist

import (
	"context"
	"log"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

const (
	pluginName = "blocklist"
)

// TODO(nikonov): domain top list

type blockListSource interface {
	IsBlocked(domain string) bool
}

// blocklistPlugin implements coredns' plugin.Handler interface
type blocklistPlugin struct {
	Next      plugin.Handler
	blockList blockListSource
}

func (b *blocklistPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	if !b.blockList.IsBlocked(state.Name()) {
		passedQueriesCount.Inc()
		return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
	}

	m := &dns.Msg{}
	m.SetReply(r)

	// non-existent domain, recursion not available
	m.Rcode = dns.RcodeNameError
	m.RecursionAvailable = false
	m.RecursionDesired = false
	// m.Answer is empty - we have nothing to reply with
	if err := w.WriteMsg(m); err != nil {
		log.Printf("ERR: failed to write answer: %v", err)
	}

	blockedQueriesCount.Inc()
	return dns.RcodeNameError, nil
}

func (*blocklistPlugin) Name() string {
	return pluginName
}

func (*blocklistPlugin) Ready() bool {
	return true
}
