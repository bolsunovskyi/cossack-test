package main

import (
	"flag"
	"log"
)

var (
	loggerHost        string
	loggerPort, speed int
)

func init() {
	flag.StringVar(&loggerHost, "lh", "127.0.0.1", "logger host")
	flag.IntVar(&loggerPort, "lp", 8989, "logger port")
	flag.IntVar(&speed, "s", 1, "number generation speed")
	flag.Parse()
}

func main() {
	prd, err := MakeProducer(speed, loggerPort, loggerHost)
	if err != nil {
		log.Fatalln(err)
	}

	if err := prd.Produce(); err != nil {
		log.Fatalln(err)
	}
}
