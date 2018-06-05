package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Logger struct {
	host                          string
	port                          int
	file                          *os.File
	buffer                        []uint64
	bufferSize, bufferPosition    int
	bufferCopies, maxBufferCopies int
	speed                         int
}

func MakeLogger(host, filePath string, port, bufferSize, bufferMaxCopies, flowSpeed int) (*Logger, error) {
	if bufferSize < 8 {
		return nil, errors.New("buffer can't be less than 8 bytes")
	}

	fp, err := os.OpenFile(filePath, os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	realBufferSize := bufferSize / 8

	return &Logger{
		port:            port,
		host:            host,
		file:            fp,
		bufferSize:      realBufferSize,
		buffer:          make([]uint64, realBufferSize),
		bufferPosition:  0,
		bufferCopies:    0,
		maxBufferCopies: bufferMaxCopies,
		speed:           flowSpeed,
	}, nil
}

func buf2uint64(buf []byte, n int) (res uint64) {
	for i := 0; i < n; i++ {
		res |= uint64(buf[i]) << uint64(i*8)
	}

	return
}

func (l *Logger) handleSocket(sock net.Conn) {
	buf := make([]byte, 8)
	throttle := time.Tick(time.Second)

	for {
		<-throttle

		for i := 0; i < l.speed; i++ {
			n, err := sock.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}

			num := buf2uint64(buf, n)
			l.buffer[l.bufferPosition] = num
			l.bufferPosition++

			if l.bufferPosition == l.bufferSize {
				if err := l.flushBuffer(); err != nil {
					log.Println(err)
					if err := sock.Close(); err != nil {
						log.Println(err)
					}
					return
				}
			}

			log.Println(num)
		}
	}
}

func (l *Logger) flushBuffer() error {
	defer func() {
		l.bufferPosition = 0
		l.buffer = make([]uint64, l.bufferSize)
	}()

	l.bufferCopies++
	if l.bufferCopies > l.maxBufferCopies {
		return errors.New("peak memory threshold")
	}

	bufferCopy := l.buffer

	go func(buff []uint64) {
		for _, n := range bufferCopy {
			l.file.WriteString(fmt.Sprintf("%d\n", n))
		}

		l.bufferCopies--
	}(bufferCopy)

	return nil
}

func (l *Logger) Listen() error {
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
