package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to decode data")
		return
	}

	block, err := generateBlock(params.Data)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to generate block")
		return
	}

	if err := blockchain.addBlock(block); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to add block")
		return
	}

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
