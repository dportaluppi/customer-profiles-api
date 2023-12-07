package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yalochat/go-commerce-components/flat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/internal/repository"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

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

	// Flattener
	f := flat.NewFlattener()

	// Attributes
	attr := profile.NewMongoRepository(ctx, mongoClient, cfg.Mongo.DB, f)

	// Entities
	router := gin.Default()
	entities := repository.NewMongoRepository[*profile.Entity](mongoClient, cfg.Mongo.DB, "entities")
	eHandler := iprofile.NewHandler(
		profile.NewSaver(entities, attr),
		profile.NewDeleter(entities, attr),
		profile.NewGetter(entities, attr),
	)

	router.POST("/accounts/:accountId/entities", eHandler.Create)
	router.PUT("/accounts/:accountId/entities/:id", eHandler.Update)
	router.DELETE("/accounts/:accountId/entities/:id", eHandler.Delete)
	router.GET("/accounts/:accountId/entities/:id", eHandler.GetByID)
	router.GET("/accounts/:accountId/entities", eHandler.GetAll)
	router.GET("/accounts/:accountId/entities/keys", eHandler.GetKeys)

	router.POST("/accounts/:accountId/entities/search", eHandler.Query)
	router.POST("/accounts/:accountId/entities/queries/jsonlogic", eHandler.QueryJsonLogic) // TODO: Remove this endpoint

	// Relationships
	router.POST("/accounts/:accountId/entities/:id/relationships", eHandler.CreateRelationship)
	router.PUT("/accounts/:accountId/entities/:id/relationships", eHandler.ReplaceRelationships)
	if err = router.Run(":8030"); err != nil {
		panic(err)
	}
}
