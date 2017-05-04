package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

const (
	DefaultConfigPath = "config.json"
)

func main() {
	var configPath = flag.String("conf", DefaultConfigPath, "path to your config file")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal(err)
	}

	newNode(config).run()
}
