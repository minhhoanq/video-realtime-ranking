package main

import (
	"video-realtime-ranking/interaction-processing-service/config"
	"video-realtime-ranking/interaction-processing-service/internal/initial"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	initial.Initial(config)
}
