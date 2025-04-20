package main

import (
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/config"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/initial"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	initial.Initial(config)
}
