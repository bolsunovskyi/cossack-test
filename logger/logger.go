package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Logger struct {
	host                          string
	port                          int
	logWriter                     ReadWriteSyncer
	buffer                        []uint64
	bufferSize, bufferPosition    int
	bufferCopies, maxBufferCopies int
	speed                         int
	clients                       []net.Conn
	listener                      net.Listener
}

type ReadWriteSyncer interface {
	io.WriteCloser
	Sync() error
}

func MakeLogger(host, filePath string, port, bufferSize, bufferMaxCopies, flowSpeed int,
	encryptLog bool, encKey string) (*Logger, error) {
	if bufferSize < 8 {
		return nil, errors.New("buffer can't be less than 8 bytes")
	}

	fp, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	realBufferSize := bufferSize / 8

	logger := Logger{
		port:            port,
		host:            host,
		logWriter:       fp,
		bufferSize:      realBufferSize,
		buffer:          make([]uint64, realBufferSize),
		bufferPosition:  0,
		bufferCopies:    0,
		maxBufferCopies: bufferMaxCopies,
		speed:           flowSpeed,
	}

	if encryptLog {
		cph, err := MakeCipher(encKey, fp)
		if err != nil {
			return nil, err
		}

		logger.logWriter = cph
	}

	return &logger, nil
}

func buf2uint64(buf []byte, n int) (res uint64) {
	for i := 0; i < n; i++ {
		res |= uint64(buf[i]) << uint64(i*8)
	}

	return
}

func (l *Logger) Stop() error {
	for k := range l.clients {
		if err := l.clients[k].Close(); err != nil {
			log.Println(err)
		}
	}

	if err := l.logWriter.Sync(); err != nil {
		log.Println(err)
	}

	if err := l.logWriter.Close(); err != nil {
		log.Println(err)
	}

	return l.listener.Close()
}

func (l *Logger) handleSocket(sock net.Conn) {
	defer func() {
		//sync file on client exit
		if err := l.flushBuffer(); err != nil {
			log.Println(err)
		}

		for i := range l.clients {
			if l.clients[i] == sock {
				l.clients = append(l.clients[:i], l.clients[i+1:]...)
			}
		}
	}()

	buf := make([]byte, 8)
	throttle := time.Tick(time.Second)

	for {
		<-throttle

		for i := 0; i < l.speed; i++ {
			n, err := sock.Read(buf)
			if err != nil {
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

	bufferCopy := l.buffer[0:l.bufferPosition]

	go func(buff []uint64) {
		for _, n := range bufferCopy {
			if _, err := l.logWriter.Write([]byte(fmt.Sprintf("%d\n", n))); err != nil {
				log.Println(err)
			}
		}
		if err := l.logWriter.Sync(); err != nil {
			log.Println(err)
		}

		l.bufferCopies--
	}(bufferCopy)

	return nil
}

func (l *Logger) Listen() (err error) {
	l.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return
	}

	go func() {
		time.Sleep(time.Second * 10)
		l.listener.Close()
	}()

	for {
		sock, err := l.listener.Accept()
		if err != nil {
			return err
		}
		l.clients = append(l.clients, sock)
		go l.handleSocket(sock)

	}
}
