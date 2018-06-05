package main

import (
	"fmt"
	"log"
	"net"
)

type Logger struct {
	host string
	port int
}

func MakeLogger(host string, port int) *Logger {
	return &Logger{port: port, host: host}
}

func buf2uint64(buf []byte, n int) (res uint64) {
	for i := 0; i < n; i++ {
		res |= uint64(buf[i]) << uint64(i*8)
	}

	return
}

func (l Logger) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)

	buf := make([]byte, 8)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(buf2uint64(buf, n))
		}
	}
}
