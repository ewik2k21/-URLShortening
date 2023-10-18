package config

import (
	"flag"
	"os"
)

var FlagA string
var FlagB string

func ParseFlags() {
	flag.StringVar(&FlagA, "a", ":8080", "port for server")
	flag.StringVar(&FlagB, "b", "localhost"+FlagA, "base address result")
	if baseUrl, err := os.LookupEnv("BASE_URL"); err {
		FlagB = baseUrl
	}
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		FlagA = envServerAddress
	}

	flag.Parse()
}
