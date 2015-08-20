package main

import (
	"log"
	"github.com/kelseyhightower/envconfig"
)

// Note: the `default` tag must appear before `envconfig` for the default thing
// to work.
type Config struct {
	Port        string `envconfig:"PORT"`
	DatabaseUrl string `envconfig:"DATABASE_URL"`
}

var config Config

func init() {
	err := envconfig.Process("sql-drain", &config)
	if err != nil {
		log.Fatalf("Incomplete config: %v", err)
	}
	log.Printf("Config => %+v", config)
}
