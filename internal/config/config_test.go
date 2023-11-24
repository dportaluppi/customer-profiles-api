package config

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yalochat/go-commerce-components/logging"

	"github.com/yalochat/go-commerce-components/config"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	getEnv := func() map[string]string {
		envs := make(map[string]string)
		for _, e := range os.Environ() {
			kv := strings.SplitN(e, "=", 2)
			envs[kv[0]] = kv[1]
		}
		return envs
	}
	initialEnv := getEnv()
	resetEnv := func(t *testing.T) {
		for k := range getEnv() {
			_, present := initialEnv[k]
			if !present {
				require.NoError(t, os.Unsetenv(k))
			}
		}
		for k, v := range initialEnv {
			require.NoError(t, os.Setenv(k, v))
		}
	}
	tests := []struct {
		it     string
		envs   func(t *testing.T) map[string]string
		assert func(t *testing.T, c *Config, err error)
	}{
		{
			it: "errors out if required env vars are missing",
			envs: func(t *testing.T) map[string]string {
				return map[string]string{}
			},
			assert: func(t *testing.T, c *Config, err error) {
				require.Nil(t, c, "config should be nil on error")
				require.Error(t, err, "error should be returned")
				require.Contains(t, err.Error(), EnvPrefix+"_TRACE_SERVICE_NAME", "not found env should return error")
			},
		},
		{
			it: "valid env vars should return hydrated config",
			envs: func(_ *testing.T) map[string]string {
				return map[string]string{
					EnvPrefix + "_TRACE_SERVICE_NAME":    "test",
					EnvPrefix + "_HEALTH_CHECK_INTERVAL": "10s",
					EnvPrefix + "_HEALTH_CHECK_TIMEOUT":  "10s",
					EnvPrefix + "_SERVER_HOST":           "test",
					EnvPrefix + "_SERVER_PORT":           "test",
					EnvPrefix + "_SERVER_READ_TIMEOUT":   "10s",
					EnvPrefix + "_SERVER_WRITE_TIMEOUT":  "10s",
					EnvPrefix + "_SERVER_METRICS_PATH":   "test",
					EnvPrefix + "_ENGINE_DEBUG":          "true",
					EnvPrefix + "_LOG_LEVEL":             "debug",
					EnvPrefix + "_LOG_FORMAT":            "text",
					EnvPrefix + "_AEROSPIKE_ADDRESS":     "127.0.0.1",
					EnvPrefix + "_AEROSPIKE_PORT":        "3000",
					EnvPrefix + "_AEROSPACE_NAMESPACE":   "test",
				}
			},
			assert: func(t *testing.T, c *Config, err error) {
				require.NoError(t, err, "no error should be returned")
				require.NotNil(t, c, "config should not be nil on success")
				require.Equal(t, &Config{
					Trace: Trace{
						ServiceName: "test",
					},
					HealthCheck: HealthCheck{
						Interval: 10 * time.Second,
						Timeout:  10 * time.Second,
					},
					Engine: Engine{
						Debug: true,
					},
					Server: Server{
						Host:         "test",
						Port:         "test",
						ReadTimeout:  10 * time.Second,
						WriteTimeout: 10 * time.Second,
						MetricsPath:  "test",
					},
					Environment: config.Environment{
						Name: "test",
					},
					Log: logging.Config{
						Level:  "debug",
						Format: "text",
					},
					Aerospike: Aerospike{
						Address:   "127.0.0.1",
						Port:      3000,
						Namespace: "customers-profiles-api",
					},
				}, c, "invalid config returned")
			},
		},
		{
			it: "default values should be populated correctly",
			envs: func(_ *testing.T) map[string]string {
				return map[string]string{
					EnvPrefix + "_TRACE_SERVICE_NAME": "test",
				}
			},
			assert: func(t *testing.T, c *Config, err error) {
				require.NoError(t, err, "no error should be returned")
				require.NotNil(t, c, "config should not be nil on success")
				require.Equal(t, &Config{
					Environment: config.Environment{
						Name: "test",
					},
					Trace: Trace{
						ServiceName: "test",
					},
					HealthCheck: HealthCheck{
						Interval: 5 * time.Second,
						Timeout:  5 * time.Second,
					},
					Engine: Engine{
						Debug: false,
					},
					Server: Server{
						Host:         "",
						Port:         "8080",
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 5 * time.Second,
						MetricsPath:  "/metrics",
					},
					Log: logging.Config{
						Level:  "info",
						Format: "json",
					},
					Aerospike: Aerospike{
						Address:   "127.0.0.1",
						Port:      3000,
						Namespace: "customers-profiles-api",
					},
				}, c, "invalid config returned")
			},
		},
	}
	for _, tt := range tests {
		for k, v := range tt.envs(t) {
			require.NoError(t, os.Setenv(k, v), "cannot set env var")
		}
		t.Run(tt.it, func(t *testing.T) {
			cfg, err := Load(context.TODO())
			tt.assert(t, cfg, err)
		})
		resetEnv(t)
	}
}
