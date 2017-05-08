package main

import (
	"errors"
	"flag"
)

var (
	apiAddr   = flag.String("api", ":3001", "api server address")
	p2pAddr   = flag.String("p2p", ":6001", "p2p server address")
	p2pOrigin = flag.String("origin", "http://127.0.0.1", "p2p origin")

	ErrInvalidChain       = errors.New("invalid chain")
	ErrInvalidBlock       = errors.New("invalid block")
	ErrUnknownMessageType = errors.New("unknown message type")
)

func main() {
	flag.Parse()
	newNode().run()
}
