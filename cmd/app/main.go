package main

import (
	"video-realtime-ranking/config"
	"video-realtime-ranking/internal/initial"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	initial.Initial(config)
}
