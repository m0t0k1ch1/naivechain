package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Blockchain []*Block

var (
	mu         = sync.RWMutex{}
	blockchain = Blockchain{genesisBlock}

	ErrInvalidBlock = errors.New("invalid block")
)

func (bc *Blockchain) addBlock(block *Block) error {
	ok, err := isValidBlock(block, blockchain.getLatestBlock())
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidBlock
	}

	mu.Lock()
	defer mu.Unlock()

	blockchain = append(blockchain, block)

	return nil
}

func (bc *Blockchain) getLatestBlock() *Block {
	mu.RLock()
	defer mu.RUnlock()

	return blockchain[len(blockchain)-1]
}

func isValidBlock(block, previousBlock *Block) (bool, error) {
	if block.Index != previousBlock.Index+1 {
		log.Println("invalid index")
		return false, nil
	}

	previousBlockHash, err := previousBlock.hash()
	if err != nil {
		return false, err
	}
	if block.PreviousHash != previousBlockHash {
		log.Println("invalid previous block hash")
		return false, nil
	}

	return true, nil
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
