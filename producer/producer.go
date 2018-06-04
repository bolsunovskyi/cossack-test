package main

import (
	"fmt"
	"time"
)

type Producer struct {
	speed int //numbers per second, 5 means 5 number for 1 second, etc...
}

func MakeProducer(speed int) *Producer {
	return &Producer{
		speed: speed,
	}
}

func (p *Producer) produce() {
	var n1, n2, n3 uint64 = 1, 1, 0
	limiter := time.Tick(time.Second)

	for {
		<-limiter

		for i := 0; i < p.speed; i++ {
			n3 = n1 + n2
			n1 = n2
			n2 = n3
			fmt.Println(n3)
		}
	}
}
