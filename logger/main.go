package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	host, logPath, encryptKey, decryptFilePath string
	port, bufferSize, bufferCopies, flowSpeed  int
	encryptFile, decryptFile                   bool
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
	flag.StringVar(&decryptFilePath, "df", "numbers.plain.log", "decrypt file path")
	flag.BoolVar(&decryptFile, "d", false, "just decrypt log file")
	flag.Parse()
}

func main() {
	if decryptFile {
		decryptLogs()
		return
	}

	logger, err := MakeLogger(host, logPath, port, bufferSize, bufferCopies, flowSpeed, encryptFile, encryptKey)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go logger.Listen()

	<-c

	logger.Stop()

	os.Exit(0)
}

func decryptLogs() {
	fp, err := os.Open(logPath)
	if err != nil {
		log.Fatalln(err)
	}

	outFile, err := os.Create(decryptFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	cph, err := MakeCipher(encryptKey, fp)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := cph.Decrypt(outFile); err != nil {
		log.Fatalln(err)
	}
}
