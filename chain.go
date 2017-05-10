package main

import (
	"sync"
	"time"
)

type Blockchain struct {
	blocks Blocks
	mu     sync.RWMutex
}

func newBlockchain() *Blockchain {
	return &Blockchain{
		blocks: Blocks{genesisBlock},
		mu:     sync.RWMutex{},
	}
}

func (bc *Blockchain) len() int {
	return len(bc.blocks)
}

func (bc *Blockchain) getGenesisBlock() *Block {
	return bc.getBlock(0)
}

func (bc *Blockchain) getLatestBlock() *Block {
	return bc.getBlock(bc.len() - 1)
}

func (bc *Blockchain) getBlock(index int) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.blocks[index]
}

func (bc *Blockchain) generateBlock(data string) *Block {
	block := &Block{
		Index:        bc.getLatestBlock().Index + 1,
		PreviousHash: bc.getLatestBlock().Hash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}
	block.Hash = block.hash()

	return block
}

func (bc *Blockchain) addBlock(block *Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.blocks = append(bc.blocks, block)
}

func (bc *Blockchain) replaceBlocks(blocks Blocks) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.blocks = blocks
}

func (bc *Blockchain) isValidGenesisBlock() bool {
	block := bc.getGenesisBlock()

	return block.Hash == genesisBlock.Hash &&
		block.isValidHash()
}

func (bc *Blockchain) isValid() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if bc.len() == 0 {
		return false
	}
	if !bc.isValidGenesisBlock() {
		return false
	}

	prevBlock := bc.getGenesisBlock()
	for i := 1; i < bc.len(); i++ {
		block := bc.getBlock(i)

		if ok := isValidBlock(block, prevBlock); !ok {
			return false
		}

		prevBlock = block
	}

	return true
}
