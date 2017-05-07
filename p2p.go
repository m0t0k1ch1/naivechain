package main

import (
	"time"

	"golang.org/x/net/websocket"
)

func (node *Node) addSocket(ws *websocket.Conn) {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.sockets = append(node.sockets, ws)
}

func (node *Node) connectToPeers(peers []string) {
	for _, peer := range peers {
		ws, err := websocket.Dial(peer, "", *p2pOrigin)
		if err != nil {
			node.logError(err)
			continue
		}
		node.addSocket(ws)

		go node.p2pHandler(ws)

		// TODO: get latest block
	}
}

func (node *Node) p2pHandler(ws *websocket.Conn) {
	for {
		// TODO: message handling
		node.log("sleeping...")
		time.Sleep(3 * time.Second)
	}
}
