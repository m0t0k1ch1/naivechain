package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

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

func (node *Node) newLatestBlockMessage() (*Message, error) {
	return newBlocksMessage(Blocks{node.blockchain.getLatestBlock()})
}

func (node *Node) newAllBlocksMessage() (*Message, error) {
	return newBlocksMessage(node.blockchain.blocks)
}

func (node *Node) broadcastLatestBlock() error {
	msg, err := node.newLatestBlockMessage()
	if err != nil {
		return err
	}

	node.broadcast(msg)

	return nil
}

func (node *Node) broadcast(msg *Message) {
	for _, conn := range node.conns {
		if err := node.send(conn, msg); err != nil {
			node.logError(err)
		}
	}
}

func (node *Node) sendLatestBlock(conn *Conn) error {
	msg, err := node.newLatestBlockMessage()
	if err != nil {
		return err
	}

	return node.send(conn, msg)
}

func (node *Node) sendAllBlocks(conn *Conn) error {
	msg, err := node.newAllBlocksMessage()
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

func (node *Node) handleBlocksResponse(conn *Conn, msg *Message) error {
	var blocks Blocks
	if err := json.Unmarshal([]byte(msg.Data), &blocks); err != nil {
		return err
	}
	sort.Sort(blocks)

	latestBlock := blocks[len(blocks)-1]
	if latestBlock.Index > node.blockchain.getLatestBlock().Index {
		if node.blockchain.getLatestBlock().Hash == latestBlock.PreviousHash {
			if isValidBlock(latestBlock, node.blockchain.getLatestBlock()) {
				node.blockchain.addBlock(latestBlock)
				if err := node.broadcastLatestBlock(); err != nil {
					return err
				}
			}
		} else if len(blocks) == 1 {
			node.broadcast(newQueryAllMessage())
		} else {
			bc := newBlockchain()
			bc.replaceBlocks(blocks)
			if bc.isValid() {
				node.blockchain.replaceBlocks(blocks)
				if err := node.broadcastLatestBlock(); err != nil {
					return err
				}
			}
		}
	}

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
		case MessageTypeQueryLatest:
			if err := node.sendLatestBlock(conn); err != nil {
				node.logError(err)
			}
		case MessageTypeQueryAll:
			if err := node.sendAllBlocks(conn); err != nil {
				node.logError(err)
			}
		case MessageTypeBlocks:
			if err := node.handleBlocksResponse(conn, &msg); err != nil {
				node.logError(err)
			}
		default:
			node.logError(ErrUnknownMessageType)
		}
	}
}
