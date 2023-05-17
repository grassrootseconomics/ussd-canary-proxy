package main

import (
	"strings"

	"github.com/grassrootseconomics/cic-custodial/pkg/logg"
	"github.com/grassrootseconomics/ussd-canary-proxy/internal/cache"
	"github.com/grassrootseconomics/ussd-canary-proxy/internal/proxy"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/zerodha/logf"
)

// Load logger.
func initLogger() logf.Logger {
	loggOpts := logg.LoggOpts{}

	if debugFlag {
		loggOpts.Color = true
		loggOpts.Caller = true
		loggOpts.Debug = true
	}

	return logg.NewLogg(loggOpts)
}

// Load config file.
func initConfig() *koanf.Koanf {
	var (
		ko = koanf.New(".")
	)

	confFile := file.Provider(confFlag)
	if err := ko.Load(confFile, toml.Parser()); err != nil {
		lo.Fatal("init: could not load config file", "error", err)
	}

	if err := ko.Load(env.Provider("CANARY_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "CANARY_")), "__", ".")
	}), nil); err != nil {
		lo.Fatal("init: could not override config from env vars", "error", err)
	}

	if debugFlag {
		ko.Print()
	}

	return ko
}

func initPgCache() cache.Cache {
	store, err := cache.NewPgCache(cache.PgCacheOpts{
		DSN:                  ko.MustString("postgres.dsn"),
		MigrationsFolderPath: migrationsFolderFlag,
		QueriesFolderPath:    queriesFlag,
	})
	if err != nil {
		lo.Fatal("init: critical error loading Postgres cache", "error", err)
	}

	return store
}

func initProxy() *proxy.Proxy {
	proxy, err := proxy.InitProxy(proxy.Opts{
		Logg:       lo,
		V1Upstream: ko.MustString("upstream.v1"),
		V2Upstream: ko.MustString("upstream.v2"),
	})
	if err != nil {
		lo.Fatal("init: critical error loading upstream proxies", "error", err)
	}

	return proxy
}
