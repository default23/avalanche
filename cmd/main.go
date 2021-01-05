package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/default23/avalanche/config"
)

var (
	configPath = ""
)

func init() {
	f := new(logrus.TextFormatter)
	f.FullTimestamp = true
	f.TimestampFormat = "02.01.2006 15:04:05"

	logrus.SetFormatter(f)

	flag.StringVar(&configPath, "config", "", "path to configuration file")
	flag.Parse()
}

func main() {
	var cfg *config.Config
	var err error

	if configPath == "" {
		cfg = config.Default()
		logrus.Warn("config path is not specified, using default configuration")
	} else {
		logrus.Infof("reading configuration from '%s'", configPath)
		cfg, err = config.Read(configPath)
		if err != nil {
			logrus.Fatalf("parse configuration file failed: %s", err)
			return
		}
	}

	logrus.Info("successfully parsed app configuration")
	logrus.Infof("instantiating the proxy service at '%s'", cfg.Server.Addr)
}
