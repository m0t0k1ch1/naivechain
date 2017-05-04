package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		// TODO: output error log
		return
	}

	block, err := generateBlock(params.Data)
	if err != nil {
		// TODO: output error log
		return
	}

	blockchain.addBlock(block)

	// TODO: bloadcast
}

func generateBlock(data string) (*Block, error) {
	previousBlock := blockchain.getLatestBlock()
	previousBlockHash, err := previousBlock.hash()
	if err != nil {
		return nil, err
	}

	return &Block{
		Index:        previousBlock.Index + 1,
		PreviousHash: previousBlockHash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}, nil
}
