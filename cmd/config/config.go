package config

import (
	"flag"
	"os"
)

var FlagPort string
var FlagBaseURL string

func ParseFlags() {
	flag.StringVar(&FlagPort, "a", ":8080", "port for server")
	flag.StringVar(&FlagBaseURL, "b", "localhost"+FlagPort, "base address result")
	if baseURL, err := os.LookupEnv("BASE_URL"); err {
		FlagBaseURL = baseURL
	}
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		FlagPort = envServerAddress
	}

	flag.Parse()
}
