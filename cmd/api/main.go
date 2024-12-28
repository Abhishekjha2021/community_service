package main

import (
	"log"

	"github.com/Abhishekjha321/community_service/cmd/api/app"
	"github.com/Abhishekjha321/community_service/pkg/config"
)

func main() {
	err := config.Initialize()
	if err != nil {
		log.Fatal("failed to initialize config: %w", err)
	}

	app := &app.Application{}
	app.Init()
	app.Start()
}
