package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/zerodha/logf"
)

type (
	Opts struct {
		Logg       logf.Logger
		V1Upstream string
		V2Upstream string
	}

	Proxy struct {
		V1 *httputil.ReverseProxy
		V2 *httputil.ReverseProxy
	}
)

func InitProxy(o Opts) (*Proxy, error) {
	o.Logg.Debug("proxy: upstream urls", "v1", o.V1Upstream, "v2", o.V2Upstream)
	v1, err := newProxy(o.V1Upstream)
	if err != nil {
		return nil, err
	}

	v2, err := newProxy(o.V2Upstream)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		V1: v1,
		V2: v2,
	}, nil
}

func newProxy(upstreamTarget string) (*httputil.ReverseProxy, error) {
	upstreamUrl, err := url.Parse(upstreamTarget)
	if err != nil {
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Header.Add("X-Forwarded-Host", r.Host)
			r.Header.Add("X-Origin-Host", upstreamUrl.Host)

			r.Host = upstreamUrl.Host

			r.URL.Scheme = upstreamUrl.Scheme
			r.URL.Host = upstreamUrl.Host
			r.URL.Path = upstreamUrl.Path
		},
	}

	return proxy, nil
}
