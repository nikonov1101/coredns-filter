package blocklist

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/pkg/errors"

	"gitlab.com/nikonov1101/coredns-filter/blocklist/piholedb"
)

func init() {
	// register blacklist as earliest plugin:
	// we don't want to bother other plugins with requests we want to drop.
	dnsserver.Directives = append([]string{pluginName}, dnsserver.Directives...)

	plugin.Register(pluginName, setup)
}

func setup(c *caddy.Controller) error {
	databasePath := "gravity.db"

	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			// use the default path
		case 1:
			databasePath = args[0]
		default:
			return errors.New("blocklist: invalid amount of args")
		}
	}

	src, err := piholedb.New(databasePath)
	if err != nil {
		return errors.Wrapf(err, "blocklist: failed to load database")
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &blocklistPlugin{Next: next, blockList: src}
	})
	return nil
}
