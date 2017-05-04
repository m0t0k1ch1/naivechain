package main

import (
	"errors"
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
		return false, nil
	}

	previousBlockHash, err := previousBlock.hash()
	if err != nil {
		return false, err
	}
	if block.PreviousHash != previousBlockHash {
		return false, nil
	}

	return true, nil
}
