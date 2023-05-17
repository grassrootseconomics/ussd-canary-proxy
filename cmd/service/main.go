package main

import (
	"flag"
	"net/http"

	"github.com/grassrootseconomics/ussd-canary-proxy/internal/api"
	"github.com/knadh/koanf/v2"
	"github.com/zerodha/logf"
)

var (
	build string

	confFlag             string
	debugFlag            bool
	migrationsFolderFlag string
	queriesFlag          string

	ko *koanf.Koanf
	lo logf.Logger
)

func init() {
	flag.StringVar(&confFlag, "config", "config.toml", "Config file location")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.StringVar(&migrationsFolderFlag, "migrations", "migrations/", "Migrations folder location")
	flag.StringVar(&queriesFlag, "queries", "queries.sql", "Queries file location")
	flag.Parse()

	lo = initLogger()
	ko = initConfig()
}

func main() {
	lo.Info("main: starting ussd-canary-proxy", "build", build)

	cache := initPgCache()
	proxy := initProxy()

	api := api.NewAPI(api.Opts{
		Cache: cache,
		Logg:  lo,
		Proxy: proxy,
	})

	router := initRouter(api)

	if ko.Bool("service.serve_proxy") {
		lo.Info("main: starting in front facing server mode")
		// CertMagic
	} else {
		lo.Info("main: starting in internal server mode")
		http.ListenAndServe(ko.MustString("service.address"), router)
	}
}
