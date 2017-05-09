package main

import (
	"encoding/json"
	"net/http"
)

func (node *Node) blocksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(node.blockchain.blocks)
	if err != nil {
		node.error(w, err, "failed to decode response")
		return
	}

	node.writeResponse(w, b)
}

func (node *Node) mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		node.error(w, err, "failed to decode params")
		return
	}

	block, err := node.blockchain.generateBlock(params.Data)
	if err != nil {
		node.error(w, err, "failed to generate block")
		return
	}

	node.blockchain.addBlock(block)

	// TODO: broadcast

	b, err := json.Marshal(map[string]string{
		"hash": block.hash(),
	})
	if err != nil {
		node.error(w, err, "failed to decode response")
	}

	node.writeResponse(w, b)
}

func (node *Node) peersHandler(w http.ResponseWriter, r *http.Request) {
	peerHosts := make([]string, len(node.conns))
	for i, conn := range node.conns {
		peerHosts[i] = conn.remoteHost()
	}

	b, err := json.Marshal(peerHosts)
	if err != nil {
		node.error(w, err, "failed to decode response")
		return
	}

	node.writeResponse(w, b)
}

func (node *Node) addPeerHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Peer string `json:"peer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		node.error(w, err, "failed to decode params")
		return
	}

	node.connectToPeers([]string{params.Peer})

	node.peersHandler(w, r)
}
