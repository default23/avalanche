package proxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"

	"github.com/default23/avalanche/config"
	"github.com/default23/avalanche/domain"
)

type handler struct {
	server     *goproxy.ProxyHttpServer
	config     config.Proxy
	clientRepo domain.ClientRepository
	logger     *logrus.Entry
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.server.ServeHTTP(w, r)
}

func Configure(
	cfg config.Proxy,
	clientRepo domain.ClientRepository,
	logger *logrus.Entry,
) http.Handler {
	h := new(handler)

	proxyHandler := goproxy.NewProxyHttpServer()
	proxyHandler.Verbose = cfg.Logging
	proxyHandler.Logger = logger

	h.server = proxyHandler
	h.config = cfg
	h.clientRepo = clientRepo
	h.logger = logger

	if cfg.Authorization.Enabled {
		h.setupAuthorization()
	}

	return h
}
