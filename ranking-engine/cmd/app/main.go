package main

import (
	"video-realtime-ranking/ranking-engine/config"
	"video-realtime-ranking/ranking-engine/internal/initial"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	initial.Initial(config)
}
