package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Note: the `default` tag must appear before `envconfig` for the default thing
// to work.
type Config struct {
	Port string `envconfig:"PORT"`
	Json bool   `envconfig:"DRAINCAT_JSON"`
}

var config Config

func init() {
	err := envconfig.Process("draincat", &config)
	if err != nil {
		log.Fatalf("Incomplete config: %v", err)
	}
}
