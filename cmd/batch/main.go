package main

import (
	"context"
	"fmt"
)

const (
	path                 = "./cmd/batch"
	profileFile          = "profiles.json"
	words                = "words.txt"
	numProfiles          = 100000
	maxAttributes        = 10
	maxChannelAttributes = 5
)

var (
	channelTypes = []string{"whatsapp", "commerce", "engagement"}
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
