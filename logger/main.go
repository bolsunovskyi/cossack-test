package main

import "flag"

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "h", "", "bind hostname")
	flag.IntVar(&port, "p", 8989, "bind port")
	flag.Parse()
}

func main() {
	logger := MakeLogger(host, port)
	logger.Listen()
}
