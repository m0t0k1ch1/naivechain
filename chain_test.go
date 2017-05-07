package main

import (
	"sync"
	"testing"
)

func newTestBlockchain(blocks []*Block) *Blockchain {
	return &Blockchain{
		blocks: blocks,
		mu:     sync.RWMutex{},
	}
}

func TestGenerateBlock(t *testing.T) {
	bc := newTestBlockchain([]*Block{genesisBlock})

	block, err := bc.generateBlock("white noise")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if block.Index != bc.getLatestBlock().Index+1 {
		t.Errorf("want %d but %d", bc.getLatestBlock().Index+1, block.Index)
	}
	if block.Data != "white noise" {
		t.Errorf("want %q but %q", "white noise", block.Data)
	}
	prevBlockHash, err := bc.getLatestBlock().hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if block.PreviousHash != prevBlockHash {
		t.Errorf("want %q but %q", prevBlockHash, block.PreviousHash)
	}
}

func TestAddBlock(t *testing.T) {
	bc := newTestBlockchain([]*Block{genesisBlock})
	block := &Block{
		Index:        1,
		PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
		Timestamp:    1494093545,
		Data:         "white noise",
	}

	if err := bc.addBlock(block); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}

	latestBlockHash, err := bc.getLatestBlock().hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	blockHash, err := block.hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != blockHash {
		t.Errorf("want %q but %q", blockHash, latestBlockHash)
	}
}

func TestReplaceBlocks(t *testing.T) {
	bc := newTestBlockchain([]*Block{genesisBlock})
	bcNew := newTestBlockchain([]*Block{
		genesisBlock,
		&Block{
			Index:        1,
			PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
			Timestamp:    1494093545,
			Data:         "white noise",
		},
	})

	if err := bc.replaceBlocks(bcNew); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}

	latestBlockHash, err := bc.getLatestBlock().hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	blockHash, err := bcNew.getBlock(1).hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != blockHash {
		t.Errorf("want %q but %q", blockHash, latestBlockHash)
	}
}
