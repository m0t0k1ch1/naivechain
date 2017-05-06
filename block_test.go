package main

import (
	"encoding/hex"
	"testing"
)

var (
	genesisBlockHex = "00000000000000000000000000000000000000000000000000000000000000000000000000000000917c5457000000006d792067656e6573697320626c6f636b2121"
)

func TestNewBlockFromBytes(t *testing.T) {
	blockBytes, err := hex.DecodeString(genesisBlockHex)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	block, err := newBlockFromBytes(blockBytes)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if block.Index != genesisBlock.Index {
		t.Errorf("want %d but %d", genesisBlock.Index, block.Index)
	}
	if block.PreviousHash != genesisBlock.PreviousHash {
		t.Errorf("want %q but %q", genesisBlock.PreviousHash, block.PreviousHash)
	}
	if block.Timestamp != genesisBlock.Timestamp {
		t.Errorf("want %d but %d", genesisBlock.Timestamp, block.Timestamp)
	}
	if block.Data != genesisBlock.Data {
		t.Errorf("want %q but %q", genesisBlock.Data, block.Data)
	}
}

func TestBlockHex(t *testing.T) {
	blockHex, err := genesisBlock.hex()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if blockHex != genesisBlockHex {
		t.Errorf("want %q but %q", genesisBlockHex, blockHex)
	}
}
