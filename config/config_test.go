package config_test

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/default23/avalanche/config"
)

func Test_Read(t *testing.T) {
	type Suite struct {
		name       string
		path       string
		setup      func(t Suite)
		teardown   func(t Suite)
		wantErr    error
		wantConfig *config.Config
	}

	tests := []Suite{
		{
			name:    "error_read_config",
			path:    "c38d7ea2-f734-44b4-bb11-e55c3473c0d6.yaml",
			wantErr: &os.PathError{Op: "open", Path: "c38d7ea2-f734-44b4-bb11-e55c3473c0d6.yaml", Err: syscall.ENOENT},
		},
		{
			name: "error_unmarshall_failed",
			path: "conf.yaml",
			setup: func(tt Suite) {
				conf := `server:addr:8080`
				err := ioutil.WriteFile(tt.path, []byte(conf), os.ModePerm)
				assert.NoError(t, err)
			},
			wantErr: &yaml.TypeError{
				Errors: []string{
					"line 1: cannot unmarshal !!str `server:...` into config.Config",
				},
			},
		},
		{
			name: "error_ssl_PEM_not_specified",
			path: "conf.yaml",
			setup: func(tt Suite) {
				c := &config.Config{
					SSL: config.SSL{
						Enabled: true,
						PemPath: "",
						KeyPath: "/etc/server.key",
					},
				}

				conf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, conf, os.ModePerm)
				assert.NoError(t, err)
			},
			wantErr: config.ErrSSLCertRequired,
		},
		{
			name: "error_ssl_KEY_not_specified",
			path: "conf.yaml",
			setup: func(tt Suite) {
				c := &config.Config{
					SSL: config.SSL{
						Enabled: true,
						PemPath: "/etc/server.pem",
						KeyPath: "",
					},
				}

				conf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, conf, os.ModePerm)
				assert.NoError(t, err)
			},
			wantErr: config.ErrSSLCertRequired,
		},
		{
			name: "error_ssl_pem_file_not_exists",
			path: "conf.yml",
			setup: func(tt Suite) {
				c := &config.Config{
					Server: config.Server{
						Addr: ":8080",
					},
					SSL: config.SSL{
						Enabled: true,
						PemPath: "cert.pem",
						KeyPath: "cert.key",
					},
				}

				yamlConf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile("cert.key", []byte("key"), os.ModePerm)
				assert.NoError(t, err)
			},
			teardown: func(tt Suite) {
				err := os.Remove("cert.key")
				assert.NoError(t, err)
			},
			wantErr: config.ErrSSLCertRequired,
		},
		{
			name: "error_ssl_key_file_not_exists",
			path: "conf.yml",
			setup: func(tt Suite) {
				c := &config.Config{
					Server: config.Server{
						Addr: ":8080",
					},
					SSL: config.SSL{
						Enabled: true,
						PemPath: "cert.pem",
						KeyPath: "cert.key",
					},
				}

				yamlConf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile("cert.pem", []byte("pem"), os.ModePerm)
				assert.NoError(t, err)
			},
			teardown: func(tt Suite) {
				err := os.Remove("cert.pem")
				assert.NoError(t, err)
			},
			wantErr: config.ErrSSLCertRequired,
		},
		{
			name: "error_passwd_path_empty",
			path: "conf.yaml",
			setup: func(tt Suite) {
				c := &config.Config{
					Proxy: config.Proxy{
						Authorization: config.Authorization{
							Enabled:  true,
							AuthPath: "",
						},
					},
				}

				yamlConf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)
			},
			wantErr: config.ErrCredentialsNotSet,
		},
		{
			name: "error_passwd_path_invalid",
			path: "conf.yaml",
			setup: func(tt Suite) {
				c := &config.Config{
					Proxy: config.Proxy{
						Authorization: config.Authorization{
							Enabled:  true,
							AuthPath: "/etc/9ff87c10-06b9-40e9-b368-d7b67a67939b.yaml",
						},
					},
				}

				yamlConf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)
			},
			wantErr: config.ErrCredentialsNotSet,
		},
		{
			name: "success",
			path: "conf.yaml",
			setup: func(tt Suite) {
				yamlConf, err := yaml.Marshal(tt.wantConfig)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.wantConfig.SSL.PemPath, []byte("pem"), os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.wantConfig.SSL.KeyPath, []byte("key"), os.ModePerm)
				assert.NoError(t, err)
			},
			teardown: func(tt Suite) {
				err := os.Remove(tt.wantConfig.SSL.PemPath)
				assert.NoError(t, err)

				err = os.Remove(tt.wantConfig.SSL.KeyPath)
				assert.NoError(t, err)
			},
			wantConfig: &config.Config{
				Server: config.Server{
					Addr: ":8080",
				},
				SSL: config.SSL{
					Enabled: true,
					PemPath: "cert.pem",
					KeyPath: "cert.key",
				},
			},
		},
		{
			name: "success_default_addr",
			path: "conf.yaml",
			setup: func(tt Suite) {
				c := *tt.wantConfig
				c.Server.Addr = ""

				yamlConf, err := yaml.Marshal(c)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.path, yamlConf, os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.wantConfig.SSL.PemPath, []byte("pem"), os.ModePerm)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.wantConfig.SSL.KeyPath, []byte("key"), os.ModePerm)
				assert.NoError(t, err)
			},
			teardown: func(tt Suite) {
				err := os.Remove(tt.wantConfig.SSL.PemPath)
				assert.NoError(t, err)

				err = os.Remove(tt.wantConfig.SSL.KeyPath)
				assert.NoError(t, err)
			},
			wantConfig: &config.Config{
				Server: config.Server{
					Addr: ":3128",
				},
				SSL: config.SSL{
					Enabled: true,
					PemPath: "cert.pem",
					KeyPath: "cert.key",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt)
			}
			if tt.teardown != nil {
				defer tt.teardown(tt)
			}

			got, err := config.Read(tt.path)

			assert.Equal(t, tt.wantConfig, got)
			assert.Equal(t, tt.wantErr, err)

			_ = os.Remove(tt.path)
		})
	}
}

func Test_Default(t *testing.T) {
	conf := config.Default()
	wantConfig := &config.Config{Server: config.Server{Addr: ":3128"}}

	assert.Equal(t, wantConfig, conf)
}
