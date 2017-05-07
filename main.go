package main

import (
	"errors"
	"flag"
)

var (
	ErrInvalidChain = errors.New("invalid chain")
	ErrInvalidBlock = errors.New("invalid block")

	apiAddr   = flag.String("api", ":3001", "api server address")
	p2pAddr   = flag.String("p2p", ":6001", "p2p server address")
	p2pOrigin = flag.String("origin", "ws://localhost", "p2p origin")
)

func main() {
	flag.Parse()
	newNode().run()
}
