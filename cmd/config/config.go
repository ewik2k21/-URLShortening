package config

import (
	"flag"
)

var FlagA string
var FlagB string

func ParseFlags() {
	flag.StringVar(&FlagA, "a", ":8080", "port for server")
	flag.StringVar(&FlagB, "b", "localhost"+FlagA, "base address result")

	flag.Parse()
}
