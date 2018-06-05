package main

import (
	"flag"
	"log"
)

var (
	host, logPath, encryptKey                 string
	port, bufferSize, bufferCopies, flowSpeed int
	encryptFile                               bool
)

func init() {
	flag.StringVar(&host, "h", "", "bind hostname")
	flag.IntVar(&port, "p", 8989, "bind port")
	flag.StringVar(&logPath, "l", "numbers.log", "log path")
	flag.IntVar(&bufferSize, "bs", 8000, "buffer size (bytes, minimum 8)")
	flag.IntVar(&bufferCopies, "bc", 100, "allowed buffer copies during file flush")
	flag.IntVar(&flowSpeed, "fs", 1000, "flow speed")
	flag.BoolVar(&encryptFile, "ef", false, "encrypt file")
	flag.StringVar(&encryptKey, "ek", "", "encrypt key")
	flag.Parse()
}

func main() {
	logger, err := MakeLogger(host, logPath, port, bufferSize, bufferCopies, flowSpeed)
	if err != nil {
		log.Fatal(err)
	}

	logger.Listen()
}
