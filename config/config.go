package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

const defaultAddr = ":3128"

var (
	ErrSSLCertRequired   = errors.New("pem and key files is required if 'ssl enable:true'")
	ErrCredentialsNotSet = errors.New("htpasswd file path should be provided if authorization enabled")
)

type validator interface {
	valid() error
}

type Config struct {
	Server `yaml:"server"`
	SSL    `yaml:"ssl"`
	Proxy  `yaml:"proxy"`
}

// valid validates the whole configuration file
func (c *Config) valid() error {
	for _, cc := range []validator{&c.SSL, &c.Server, &c.Proxy.Authorization} {
		if err := cc.valid(); err != nil {
			return err
		}
	}

	return nil
}

type Proxy struct {
	Logging       bool `yaml:"logging"`
	Authorization `yaml:"authorization"`
}

type Authorization struct {
	Enabled  bool   `yaml:"enabled"`
	AuthPath string `yaml:"passwdfile"`
}

// valid validates the Authorization configuration
// if auth is enabled, path to passwd file
// should be provided (AuthPath) and exists in FS
func (ac *Authorization) valid() error {
	if ac.Enabled && ac.AuthPath == "" {
		return ErrCredentialsNotSet
	}

	if ac.Enabled {
		_, err := ioutil.ReadFile(ac.AuthPath)
		if err != nil {
			logrus.Errorf("unable to read the auth file from '%s': %s", ac.AuthPath, err)
			return ErrCredentialsNotSet
		}
	}

	return nil
}

type Server struct {
	Addr string `yaml:"addr"`
}

// valid validates the Server configuration
// if Addr is not provided, sets the default port:`defaultAddr`
func (sc *Server) valid() error {
	if sc.Addr == "" {
		logrus.Warnf("server address is not specified, '%s' will be used as default", defaultAddr)
		sc.Addr = defaultAddr
	}

	return nil
}

type SSL struct {
	Enabled bool   `yaml:"enabled"`
	PemPath string `yaml:"pem"`
	KeyPath string `yaml:"key"`
}

// valid validates the SSL configuration
// if SSL enabled, the key and pem paths should be provided
func (sc *SSL) valid() error {
	if sc.Enabled {
		if sc.KeyPath == "" || sc.PemPath == "" {
			return ErrSSLCertRequired
		}

		// check the cert files exists in FS
		for _, f := range []string{sc.KeyPath, sc.PemPath} {
			_, err := ioutil.ReadFile(f)
			if err != nil {
				logrus.Errorf("unable to read the cert file from '%s': %s", f, err)
				return ErrSSLCertRequired
			}
		}
	}

	return nil
}

// Default creates an basic application
// config, without any options
func Default() *Config {
	return &Config{Server: Server{Addr: defaultAddr}}
}

// Read parses the yml config at path
func Read(path string) (*Config, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	out := &Config{}
	err = yaml.Unmarshal(conf, out)
	if err != nil {
		return nil, err
	}

	if err = out.valid(); err != nil {
		return nil, err
	}

	return out, nil
}
