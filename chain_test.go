package main

import "testing"

var testchain = Blockchain{genesisBlock}

func TestMineBlock(t *testing.T) {
	prevBlock := testchain.getLatestBlock()
	prevBlockHash, err := prevBlock.hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	data := "white noise"
	block, err := testchain.generateBlock(data)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if block.Index != prevBlock.Index+1 {
		t.Errorf("want %d but %d", prevBlock.Index+1, block.Index)
	}
	if block.Data != data {
		t.Errorf("want %q but %q", data, block.Data)
	}
	if block.PreviousHash != prevBlockHash {
		t.Errorf("want %q but %q", prevBlockHash, block.PreviousHash)
	}

	if err := testchain.addBlock(block); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	latestBlock := testchain.getLatestBlock()
	latestBlockHash, err := latestBlock.hash()
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
