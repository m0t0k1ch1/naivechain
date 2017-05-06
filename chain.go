package main

import (
	"errors"
	"sync"
	"time"
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

func (bc *Blockchain) generateBlock(data string) (*Block, error) {
	prevBlock := blockchain.getLatestBlock()
	prevBlockHash, err := prevBlock.hash()
	if err != nil {
		return nil, err
	}

	return &Block{
		Index:        prevBlock.Index + 1,
		PreviousHash: prevBlockHash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}, nil
}

func isValidBlock(block, prevBlock *Block) (bool, error) {
	if block.Index != prevBlock.Index+1 {
		return false, nil
	}

	prevBlockHash, err := prevBlock.hash()
	if err != nil {
		return false, err
	}
	if block.PreviousHash != prevBlockHash {
		return false, nil
	}

	return true, nil
}
