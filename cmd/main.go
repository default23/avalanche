package main

import (
	"flag"
	"net/http"

	"github.com/default23/avalanche/config"
	"github.com/default23/avalanche/logging"
	"github.com/default23/avalanche/proxy"
)

var (
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to configuration file")
	flag.Parse()
}

func main() {
	var cfg *config.Config
	var err error

	logger := logging.New()

	if configPath == "" {
		cfg = config.Default()
		logger.Warn("config path is not specified, using default configuration")
	} else {
		logger.Infof("reading configuration from '%s'", configPath)
		cfg, err = config.Read(configPath)
		if err != nil {
			logger.Fatalf("parse configuration file failed: %s", err)
			return
		}
	}

	logger.Info("successfully parsed app configuration")
	logger.Info("instantiating the proxy service")

	handler := proxy.Configure(true, logger)

	if cfg.SSL.Enabled {
		logger.Infof("starting the TLS server, available at %s", cfg.Server.Addr)
		logger.Fatal(http.ListenAndServeTLS(cfg.Server.Addr, cfg.SSL.PemPath, cfg.SSL.KeyPath, handler))
	} else {
		logger.Infof("starting the HTTP server, available at %s", cfg.Server.Addr)
		logger.Fatal(http.ListenAndServe(cfg.Server.Addr, handler))

	}
}
