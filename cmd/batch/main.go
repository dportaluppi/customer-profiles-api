package main

import (
	"context"
	"fmt"
)

const (
	path          = "./cmd/batch"
	profileFile   = "profiles.json"
	words         = "words.txt"
	numProfiles   = 100000
	maxAttributes = 10
)

var (
	entityTypes        = []string{"Contact", "Store"}
	relationshipsTypes = []string{"buys_from", "buys_for", "sells_for"}
)

func main() {
	profilePath := fmt.Sprintf("%s/%s", path, profileFile)
	wordsPath := fmt.Sprintf("%s/%s", path, words)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := mockProfiles(ctx, profilePath, wordsPath); err != nil {
		panic(err)
	}

	if err := insertProfiles(ctx, profilePath); err != nil {
		panic(err)
	}
}
