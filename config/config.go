package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const defaultAddr = ":3128"

var (
	ErrSSLCertRequired = errors.New("pem and key files is required if 'ssl enable:true'")
)

type Config struct {
	Server `yaml:"server"`
	SSL    `yaml:"ssl"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type SSL struct {
	Enabled bool   `yaml:"enabled"`
	PemPath string `yaml:"pem"`
	KeyPath string `yaml:"key"`
}

func (c *Config) valid() error {
	if c.SSL.Enabled {
		if c.SSL.KeyPath == "" || c.SSL.PemPath == "" {
			return ErrSSLCertRequired
		}

		// check the cert files exists in FS
		for _, f := range []string{c.SSL.KeyPath, c.SSL.PemPath} {
			_, err := ioutil.ReadFile(f)
			if err != nil {
				return ErrSSLCertRequired
			}
		}
	}

	if c.Server.Addr == "" {
		c.Server.Addr = defaultAddr
	}

	return nil
}

func Default() *Config {
	return &Config{
		Server: Server{Addr: defaultAddr},
		SSL:    SSL{Enabled: false},
	}
}

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
