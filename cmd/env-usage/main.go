package main

import (
	"github.com/dportaluppi/customer-profiles-api/internal/config"
	"github.com/yalochat/go-commerce-components/configs"
)

func main() {
	configs.UsageMain[config.Config](config.EnvPrefix)
}
