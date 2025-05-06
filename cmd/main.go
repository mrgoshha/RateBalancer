package main

import (
	"RateBalancer/internal/app"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "Path to configuration file")

	flag.Parse()
	a, err := app.NewApp(*configPath)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	a.Run()
}
