package config

import (
	"flag"
)

var FlagA string
var FlagB string

func ParseFlags() {
	flag.StringVar(&FlagA, "a", ":8888", "port for server")
	flag.StringVar(&FlagB, "b", "http://localhost", "base address result")
	flag.Parse()
}
