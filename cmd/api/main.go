package main

import (
	"context"
	"github.com/aerospike/aerospike-client-go/v6"
	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	// Aerospike
	// create an aerospike client
	client, err := aerospike.NewClient(cfg.Aerospike.Address, cfg.Aerospike.Port)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Profile
	repo := iprofile.NewAerospikeRepository(client, cfg.Aerospike.Namespace)
	profileHandler := iprofile.NewHandler(
		profile.NewUpserter(repo),
		profile.NewDeleter(repo),
		profile.NewGetter(repo),
	)

	// Routes
	router := gin.Default()
	router.POST("/profiles", profileHandler.Create)
	router.PUT("/profiles/:id", profileHandler.Update)
	router.DELETE("/profiles/:id", profileHandler.Delete)
	router.GET("/profiles/:id", profileHandler.GetByID)
	router.GET("/profiles", profileHandler.GetAll)

	if err = router.Run(":8030"); err != nil {
		panic(err)
	}
}
