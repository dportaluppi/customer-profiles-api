package config

import (
	"context"
	"time"

	"github.com/yalochat/go-commerce-components/logging"

	"github.com/yalochat/go-commerce-components/config"
)

const EnvPrefix = "UCP"

type Trace struct {
	ServiceName string `required:"true" split_words:"true"`
}

type HealthCheck struct {
	Interval time.Duration `default:"5s"`
	Timeout  time.Duration `default:"5s"`
}

type Engine struct {
	Debug bool `default:"false"`
}

type Server struct {
	Host         string        `default:""`
	Port         string        `default:"8080"`
	ReadTimeout  time.Duration `split_words:"true" default:"5s"`
	WriteTimeout time.Duration `split_words:"true" default:"5s"`
	MetricsPath  string        `split_words:"true" default:"/metrics"`
}

type Aerospike struct {
	Address   string `default:"127.0.0.1"`
	Port      int    `default:"3000"`
	Namespace string `default:"customers-profiles-api"`
}

type Mongo struct {
	Uri               string        `required:"true" default:"mongodb://localhost:27017"`
	Timeout           time.Duration `default:"1s"`
	ConnectionTimeout time.Duration `split_words:"true" default:"1s"`
	DB                string        `default:"customers-profiles-api"`
}

type Config struct {
	Environment config.Environment
	Trace       Trace
	HealthCheck HealthCheck `split_words:"true"`
	Engine      Engine
	Server      Server
	Log         logging.Config
	Aerospike   Aerospike
	Mongo       Mongo
}

// Load returns a hydrated Config object for the current environment.
// Returns an error if env vars have invalid values (non-int string for numbers,
// non-bool for booleans, etc.) or if required values are missing.
func Load(ctx context.Context) (*Config, error) {
	return config.Load[Config](ctx, EnvPrefix)
}
