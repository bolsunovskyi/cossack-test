package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Producer struct {
	loggerHost string
	loggerPort int
	conn       *net.UDPConn
	speed      int //numbers per second, 5 means 5 number for 1 second, etc...

}

func MakeProducer(speed, loggerPort int, loggerHost string) (*Producer, error) {
	servAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", loggerHost, loggerPort))
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", localAddr, servAddr)
	if err != nil {
		return nil, err
	}

	return &Producer{
		speed: speed,
		conn:  conn,
	}, nil
}

func (p Producer) sendNum(num uint64) error {
	buf := make([]byte, 8)
	var mask uint64 = 0xFF
	for i := 0; i < 8; i++ {
		buf[i] = byte((num & (mask << uint64(i*8))) >> uint64(i*8))
	}

	if _, err := p.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func (p Producer) Produce() error {
	var n1, n2, n3 uint64 = 1, 1, 0
	limiter := time.Tick(time.Second)

	for {
		<-limiter

		for i := 0; i < p.speed; i++ {
			n3 = n1 + n2
			n1 = n2
			if err := p.sendNum(n3); err != nil {
				log.Println(err)
				n1 = n3 - n2
				continue
			}
			fmt.Println(n3)
			n2 = n3
		}
	}
}
