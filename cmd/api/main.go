package main

import (
	"context"
	"log"

	"github.com/aerospike/aerospike-client-go/v6"
	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/internal/repository"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	// Aerospike client
	client, err := aerospike.NewClient(cfg.Aerospike.Address, cfg.Aerospike.Port)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Set up MongoDB client
	clientOptions := options.Client().
		ApplyURI(cfg.Mongo.Uri).
		SetConnectTimeout(cfg.Mongo.ConnectionTimeout).
		SetSocketTimeout(cfg.Mongo.Timeout)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Entities
	router := gin.Default()
	entities := repository.NewMongoRepository[*profile.Entity](mongoClient, cfg.Mongo.DB, "entities")
	eHandler := iprofile.NewHandler(
		profile.NewSaver(entities),
		profile.NewDeleter(entities),
		profile.NewGetter(entities),
	)
	router.POST("/accounts/:accountId/entities", eHandler.Create)
	router.PUT("/accounts/:accountId/entities/:id", eHandler.Update)
	router.DELETE("/accounts/:accountId/entities/:id", eHandler.Delete)
	router.GET("/accounts/:accountId/entities/:id", eHandler.GetByID)
	router.GET("/accounts/:accountId/entities", eHandler.GetAll)

	router.POST("/accounts/:accountId/entities/search", eHandler.Query)
	router.POST("/accounts/:accountId/entities/queries/jsonlogic", eHandler.QueryJsonLogic) // TODO: Remove this endpoint

	// Relationships
	router.POST("/accounts/:accountId/entities/:id/relationships", eHandler.CreateRelationship)
	router.PUT("/accounts/:accountId/entities/:id/relationships", eHandler.ReplaceRelationships)
	if err = router.Run(":8030"); err != nil {
		panic(err)
	}
}
