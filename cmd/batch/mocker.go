package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
)

func mockProfiles(_ context.Context, profileFile, wordsFile string) error {
	dictionary, err := loadDictionary(wordsFile)
	if err != nil {
		return errors.New("failed loading dictionary")
	}

	f, err := CreateFile(profileFile)
	if err != nil {
		return errors.New("failed creating profiles file")
	}

	if err = f.Start(); err != nil {
		return errors.New("failed starting profiles file")
	}

	for i := 0; i < numProfiles-1; i++ {
		if err = f.Append(createProfile(dictionary)); err != nil {
			return errors.New("failed starting profiles file")
		}
	}

	if err = f.Finish(createProfile(dictionary)); err != nil {
		return errors.New("failed finished profiles file")
	}

	fmt.Println("profiles.json created successfully")

	return nil
}

func createProfile(dictionary []string) profile.Entity {
	now := time.Now()
	p := profile.Entity{
		Attributes:    attributes(maxAttributes, dictionary),
		AccountID:     "1",
		Type:          entityType(),
		Relationships: relationships(),
		CreatedAt:     &now,
	}

	if !isNewProfile() {
		r := rand.Intn(3600)
		u := now.Add(time.Duration(r) * time.Minute)

		p.UpdatedAt = &u
	}

	return p
}

func isNewProfile() bool {
	return rand.Intn(2) == 0
}

func entityType() string {
	return entityTypes[rand.Intn(len(entityTypes))]
}

func relationships() []profile.Relationship {
	var result []profile.Relationship

	for j := 0; j < rand.Intn(len(relationshipsTypes))+1; j++ {
		result = append(result, profile.Relationship{
			Type:     relationshipsTypes[j],
			TargetID: primitive.NewObjectID().Hex(),
		})
	}

	return result
}

func attributes(length int, keys []string) map[string]any {
	result := make(map[string]any)
	for j := 0; j < rand.Intn(length)+1; j++ {
		result[randomWord(keys)] = randomWord(keys)
	}

	return result
}

func randomWord(dictionary []string) string {
	return dictionary[rand.Intn(len(dictionary))]
}

func loadDictionary(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
