package config

import (
	"flag"
	"os"
)

var (
	FlagPort             string
	FlagBaseURL          string
	FlagLogLevel         string
	FlagFileName         string
	FlagConnectionString string
)

func ParseFlags() {
	flag.StringVar(&FlagPort, "a", ":8080", "port for server")
	flag.StringVar(&FlagBaseURL, "b", "localhost"+FlagPort, "base address result")
	flag.StringVar(&FlagLogLevel, "l", "Info", "logger level")
	flag.StringVar(&FlagFileName, "f", "/tmp/short-url-db.json", "file name for json")
	flag.StringVar(&FlagConnectionString, "d", "host=localhost port=5432 user=postgres password=zaxsaqswq1w2 dbname=yandex", "connection string")
	if baseURL, err := os.LookupEnv("BASE_URL"); err {
		FlagBaseURL = baseURL
	}
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		FlagPort = envServerAddress
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
	if envFileName := os.Getenv("FILE_STORAGE_PATH"); envFileName != "" {
		FlagFileName = envFileName
	}
	if envConnectionString := os.Getenv("DATABASE_DSN"); envConnectionString != "" {
		FlagConnectionString = envConnectionString
	}
	flag.Parse()
}
