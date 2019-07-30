/*
Package main is a secret service entry-point.

The secres servise helps to store secrets and get them by unique address.
*/
package main

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ilyakaznacheev/secret/internal/config"
	"github.com/ilyakaznacheev/secret/internal/server"
)

func main() {
	var conf config.Config

	// process flags and update help function
	flag.Usage = cleanenv.Usage(&conf, nil, flag.Usage)
	flag.Parse()

	// read config
	cleanenv.ReadEnv(&conf)

	// Run service
	if err := server.Run(conf); err != nil {
		log.Fatal(err)
	}
}
