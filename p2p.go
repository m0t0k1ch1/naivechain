package main

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

func (node *Node) addConn(conn *Conn) {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.conns = append(node.conns, conn)
}

func (node *Node) deleteConn(id int64) {
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

func (node *Node) disconnectPeer(conn *Conn) {
	defer conn.Close()
	node.log("disconnect peer:", conn.remoteHost())
	node.deleteConn(conn.id)
}

func (node *Node) queryLatest(conn *Conn) error {
	return node.send(conn, newQueryLatestMessage())
}

func (node *Node) queryAll(conn *Conn) error {
	return node.send(conn, newQueryAllMessage())
}

func (node *Node) responseLatest(conn *Conn) error {
	msg, err := newBlocksMessage([]*Block{node.blockchain.getLatestBlock()})
	if err != nil {
		return err
	}

	return node.send(conn, msg)
}

func (node *Node) responseAll(conn *Conn) error {
	msg, err := newBlocksMessage(node.blockchain.blocks)
	if err != nil {
		return err
	}

	return node.send(conn, msg)
}

func (node *Node) send(conn *Conn, msg *Message) error {
	node.log(fmt.Sprintf(
		"send %s message to %s",
		msg.Type.name(), conn.remoteHost(),
	))

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return websocket.Message.Send(conn.Conn, b)
}

// TODO
func (node *Node) handleBlocksResponse(msg Message) error {
	return nil
}

func (node *Node) p2pHandler(conn *Conn) {
	for {
		var b []byte
		if err := websocket.Message.Receive(conn.Conn, &b); err != nil {
			if err == io.EOF {
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

		node.log(fmt.Sprintf(
			"receive %s message from %s",
			msg.Type.name(), conn.remoteHost(),
		))

		switch msg.Type {
		case messageTypeQueryLatest:
			if err := node.responseLatest(conn); err != nil {
				node.logError(err)
			}
		case messageTypeQueryAll:
			if err := node.responseAll(conn); err != nil {
				node.logError(err)
			}
		case messageTypeResponseBlocks:
			if err := node.handleBlocksResponse(msg); err != nil {
				node.logError(err)
			}
		default:
			node.logError(ErrUnknownMessageType)
		}
	}
}
