package main

import (
	"encoding/json"
	"io"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

const (
	messageTypeQueryLatest = iota
	messageTypeQueryAll
	messageTypeResponseBlockchain
)

type Message struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type Conn struct {
	*websocket.Conn
	id int64
}

func newConn(ws *websocket.Conn) *Conn {
	return &Conn{
		Conn: ws,
		id:   time.Now().UnixNano(),
	}
}

func (conn *Conn) remoteHost() string {
	u, _ := url.Parse(conn.RemoteAddr().String())

	return u.Host
}

func (node *Node) addConn(conn *Conn) {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.conns = append(node.conns, conn)
}

func (node *Node) deleteConnection(id int64) {
	node.mu.Lock()
	defer node.mu.Unlock()

	conns := []*Conn{}
	for _, conn := range node.conns {
		if conn.id != id {
			conns = append(conns, conn)
		}
	}

	node.conns = conns
}

func (node *Node) connectToPeers(peers []string) {
	for _, peer := range peers {
		ws, err := websocket.Dial(peer, "", *p2pOrigin)
		if err != nil {
			node.logError(err)
			continue
		}

		conn := newConn(ws)
		node.addConn(conn)
		go node.p2pHandler(conn)

		// TODO: get latest block
	}
}

func (node *Node) disconnectPeer(conn *Conn) {
	defer conn.Close()
	node.deleteConnection(conn.id)
}

func (node *Node) p2pHandler(conn *Conn) {
	for {
		var b []byte
		if err := websocket.Message.Receive(conn.Conn, &b); err != nil {
			if err == io.EOF {
				node.log("disconnect peer:", conn.remoteHost())
				node.disconnectPeer(conn)
				break
			}
			node.logError(err)
			continue
		}

		var msg Message
		if err := json.Unmarshal(b, &msg); err != nil {
			node.logError(err)
			continue
		}

		switch msg.Type {
		case messageTypeQueryLatest:
			// TODO
		case messageTypeQueryAll:
			// TODO
		case messageTypeResponseBlockchain:
			// TODO
		default:
			node.logError(ErrUnknownMessageType)
		}
	}
}
