package main

import (
	"encoding/hex"
	"testing"
)

const (
	genesisBlockHex = "00000000000000000000000000000000000000000000000000000000000000000000000000000000917c5457000000006d792067656e6573697320626c6f636b2121"
)

func TestNewBlockFromBytes(t *testing.T) {
	blockBytes, err := hex.DecodeString(genesisBlockHex)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	block, err := NewBlockFromBytes(blockBytes)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if block.Index != GenesisBlock.Index {
		t.Errorf("want %d but %d", GenesisBlock.Index, block.Index)
	}
	if block.PreviousHash != GenesisBlock.PreviousHash {
		t.Errorf("want %q but %q", GenesisBlock.PreviousHash, block.PreviousHash)
	}
	if block.Timestamp != GenesisBlock.Timestamp {
		t.Errorf("want %d but %d", GenesisBlock.Timestamp, block.Timestamp)
	}
	if block.Data != GenesisBlock.Data {
		t.Errorf("want %q but %q", GenesisBlock.Data, block.Data)
	}
}

func TestBlockHex(t *testing.T) {
	blockHex, err := GenesisBlock.Hex()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if blockHex != genesisBlockHex {
		t.Errorf("want %q but %q", genesisBlockHex, blockHex)
	}
}
