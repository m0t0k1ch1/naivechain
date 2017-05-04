package main

import "sync"

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
