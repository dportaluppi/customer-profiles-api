package main

import (
	"context"
	"github.com/aerospike/aerospike-client-go/v6"
	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/internal/repository"
	istore "github.com/dportaluppi/customer-profiles-api/internal/store"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
	"github.com/dportaluppi/customer-profiles-api/pkg/store"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

	// Profile
	// repo := iprofile.NewAerospikeRepository(client, cfg.Aerospike.Namespace) // TODO: Remove this line and use only mongo.
	// repo = iprofile.NewMongoRepository(mongoClient, cfg.Mongo.DB)
	users := repository.NewMongoRepository[*profile.Profile](mongoClient, cfg.Mongo.DB, "users")
	uHandler := iprofile.NewHandler(
		profile.NewSaver(users),
		profile.NewDeleter(users),
		profile.NewGetter(users),
	)

	// Profiles
	router := gin.Default()
	router.POST("/profiles", uHandler.Create)
	router.PUT("/profiles/:id", uHandler.Update)
	router.DELETE("/profiles/:id", uHandler.Delete)
	router.GET("/profiles/:id", uHandler.GetByID)
	router.GET("/profiles", uHandler.GetAll)
	router.POST("/profiles/query", uHandler.Query)
	router.POST("/profiles/queries/jsonlogic", uHandler.QueryJsonLogic)

	// Stores
	stores := repository.NewMongoRepository[*store.Store](mongoClient, cfg.Mongo.DB, "stores")
	sHandler := istore.NewHandler(
		store.NewSaver(stores),
		store.NewDeleter(stores),
		store.NewGetter(stores),
	)
	router.POST("/stores", sHandler.Create)
	router.PUT("/stores/:id", sHandler.Update)
	router.DELETE("/stores/:id", sHandler.Delete)
	router.GET("/stores/:id", sHandler.GetByID)
	router.GET("/stores", sHandler.GetAll)
	router.POST("/stores/query", sHandler.Query)
	router.POST("/stores/queries/jsonlogic", sHandler.QueryJsonLogic)

	if err = router.Run(":8030"); err != nil {
		panic(err)
	}
}
