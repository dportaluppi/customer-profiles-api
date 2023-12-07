package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/yalochat/go-commerce-components/flat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dportaluppi/customer-profiles-api/internal/config"
	iprofile "github.com/dportaluppi/customer-profiles-api/internal/profile"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
)

func insertProfiles(ctx context.Context, profileFile string) error {
	pCh, err := loadProfiles(profileFile)
	if err != nil {
		return errors.New("error loading profile file")
	}

	s := upserter(ctx)

	for p := range pCh {
		if _, err = s.Create(ctx, &p); err != nil {
			return err
		}
	}

	fmt.Println("profiles inserted successfully")

	return nil
}

func loadProfiles(profileFile string) (<-chan profile.Profile, error) {
	profilesChannel := make(chan profile.Profile)

	f, err := os.Open(profileFile)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(profilesChannel)
		defer f.Close()

		decoder := json.NewDecoder(f)
		_, err = decoder.Token()
		if err != nil {
			fmt.Println("Error reading array start token:", err)
			return
		}

		for decoder.More() {
			var p profile.Profile
			if err := decoder.Decode(&p); err != nil {
				fmt.Println("Error decoding profile:", err)
				return
			}

			profilesChannel <- p
		}
	}()

	return profilesChannel, nil
}

func upserter(ctx context.Context) profile.Upserter {
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

	// Profile
	repo := iprofile.NewMongoRepository(mongoClient, cfg.Mongo.DB)

	// Flattener
	f := flat.NewFlattener()

	// Attributes
	attr := profile.NewMongoRepository(ctx, mongoClient, cfg.Mongo.DB, f)

	return profile.NewUpserter(repo, attr)
}
