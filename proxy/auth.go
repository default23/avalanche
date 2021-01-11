package proxy

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

const authHeader = "Proxy-Authorization"

func reject(r *http.Request, msg string) *http.Response {
	return goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusUnauthorized, msg)
}

func (h *handler) setupAuthorization() {
	h.server.
		OnRequest().
		DoFunc(
			func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				logger := h.logger.WithField("remote", req.RemoteAddr)

				credentials, res := h.extractCredentials(req)
				if res != nil {
					return req, res
				}

				logger = logger.WithField("user", credentials[0])
				logger.Info("client identified, validating the client and password")

				client, err := h.clientRepo.Read(credentials[0])
				if err != nil {
					logger.Errorf("failed to find the client credentials in repository: %s", err)
					return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, "")
				}
				if client == nil {
					logger.Errorf("user not found")
					return req, reject(req, "user or password not match")
				}

				if ok := client.ComparePassword(credentials[1]); !ok {
					logger.Errorf("provided password not match")
					return req, reject(req, "user or password not match")
				}

				return req, nil
			},
		)
}

func (h *handler) extractCredentials(req *http.Request) ([]string, *http.Response) {
	logger := h.logger.WithField("remote", req.RemoteAddr)
	logger.Info("authorizing the client")

	auth := req.Header.Get(authHeader)
	if auth == "" {
		return nil, reject(req, "authorization credentials is not provided")
	}

	encoded := strings.Replace(auth, "Basic ", "", -1)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		logger.Errorf("failed to base64.decode the authorization header: %s", auth)
		return nil, reject(req, "can't decode the authorization header. Service using the Basic Authorization")
	}

	credentials := strings.Split(string(decoded), ":")
	if len(credentials) != 2 {
		logger.Errorf("client provided wrong credentials data")
		return nil, reject(req, "wrong authorization credentials")
	}

	return credentials, nil
}
