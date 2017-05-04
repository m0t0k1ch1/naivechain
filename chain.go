package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Blockchain []*Block

var (
	mu         = sync.RWMutex{}
	blockchain = Blockchain{genesisBlock}
)

func (bc *Blockchain) addBlock(block *Block) {
	// TODO: validation
	mu.Lock()
	defer mu.Unlock()

	blockchain = append(blockchain, block)
}

func (bc *Blockchain) getLatestBlock() *Block {
	mu.RLock()
	defer mu.RUnlock()

	return blockchain[len(blockchain)-1]
}

func blocksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(blockchain)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to decode blockchain")
		return
	}

	w.Write(b)
}
