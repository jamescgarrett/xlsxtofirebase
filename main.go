package main

import (
	"flag"
	"log"
	xtb "xlsxtofirebase/xlsxtofirebase"
)

func main() {
	flagConfig := flag.String("config", "../config.yaml", "Path to config")
	flag.Parse()

	config, err := xtb.SetupConfig(*flagConfig)
	if err != nil {
		log.Fatalf("config: %v\n", err)
	}

	firebase, err := xtb.SetupFirebase(config)
	if err != nil {
		log.Fatalf("firebase: %v\n", err)
	}

	err = firebase.SeedDatabase(config)
	if err != nil {
		log.Fatalf("firebase: %v\n", err)
	}
}
