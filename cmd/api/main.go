package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	// Aerospike client
	/*
		client, err := aerospike.NewClient(cfg.Aerospike.Address, cfg.Aerospike.Port)
		if err != nil {
			panic(err)
		}
		defer client.Close()
	*/

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

	// Attributes
	attr := profile.NewMongoRepository(mongoClient, cfg.Mongo.DB)

	// Profile
	//repo := iprofile.NewAerospikeRepository(client, cfg.Aerospike.Namespace) // TODO: Remove this line and use only mongo.
	repo := iprofile.NewMongoRepository(mongoClient, cfg.Mongo.DB)
	profileHandler := iprofile.NewHandler(
		profile.NewUpserter(repo, attr),
		profile.NewDeleter(repo, attr),
		profile.NewGetter(repo),
	)

	// Routes
	router := gin.Default()
	router.POST("/profiles", profileHandler.Create)
	router.PUT("/profiles/:id", profileHandler.Update)
	router.DELETE("/profiles/:id", profileHandler.Delete)
	router.GET("/profiles/:id", profileHandler.GetByID)
	router.GET("/profiles", profileHandler.GetAll)
	router.GET("/profiles/keys", profileHandler.GetKeys)

	router.POST("/profiles/query", profileHandler.Query)

	if err = router.Run(":8030"); err != nil {
		panic(err)
	}
}
