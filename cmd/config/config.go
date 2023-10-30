package config

import (
	"flag"
	"os"
)

var (
	FlagPort     string
	FlagBaseURL  string
	FlagLogLevel string
)

func ParseFlags() {
	flag.StringVar(&FlagPort, "a", ":8080", "port for server")
	flag.StringVar(&FlagBaseURL, "b", "localhost"+FlagPort, "base address result")
	flag.StringVar(&FlagLogLevel, "l", "Info", "logger level")
	if baseURL, err := os.LookupEnv("BASE_URL"); err {
		FlagBaseURL = baseURL
	}
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		FlagPort = envServerAddress
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}

	flag.Parse()
}
