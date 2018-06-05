package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Logger struct {
	host       string
	port       int
	file       *os.File
	buffers    [][]byte
	bufferSize int
}

func MakeLogger(host, filePath string, port, bufferSize int) (*Logger, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &Logger{
		port:       port,
		host:       host,
		file:       fp,
		bufferSize: bufferSize,
	}, nil
}

func buf2uint64(buf []byte, n int) (res uint64) {
	for i := 0; i < n; i++ {
		res |= uint64(buf[i]) << uint64(i*8)
	}

	return
}

func (l Logger) handleSocket(sock net.Conn) {
	buf := make([]byte, 8)

	for {
		n, err := sock.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(buf2uint64(buf, n))
	}
}

func (l Logger) Listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return err
	}

	for {
		sock, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go l.handleSocket(sock)
	}
}
