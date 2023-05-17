package api

import (
	"github.com/grassrootseconomics/ussd-canary-proxy/internal/cache"
	"github.com/grassrootseconomics/ussd-canary-proxy/internal/proxy"
	"github.com/zerodha/logf"
)

type (
	Opts struct {
		Cache cache.Cache
		Logg  logf.Logger
		Proxy *proxy.Proxy
	}

	API struct {
		cache cache.Cache
		logg  logf.Logger
		proxy *proxy.Proxy
	}
)

func NewAPI(o Opts) *API {
	return &API{
		cache: o.Cache,
		logg:  o.Logg,
		proxy: o.Proxy,
	}
}
