package proxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"
)

func Configure(verbose bool, logger *logrus.Entry) http.Handler {
	proxyHandler := goproxy.NewProxyHttpServer()
	proxyHandler.Verbose = verbose
	proxyHandler.Logger = logger

	return proxyHandler
}
