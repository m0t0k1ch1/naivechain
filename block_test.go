package main

import (
	"encoding/hex"
	"testing"
)

func TestNewBlockFromBytes(t *testing.T) {
	blockBytes, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000917c5457000000006d792067656e6573697320626c6f636b2121")
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
	if blockHex != "00000000000000000000000000000000000000000000000000000000000000000000000000000000917c5457000000006d792067656e6573697320626c6f636b2121" {
		t.Errorf("want \"00000000000000000000000000000000000000000000000000000000000000000000000000000000917c5457000000006d792067656e6573697320626c6f636b2121\" but %q", blockHex)
	}
}

type isValidBlockTestCase struct {
	block     *Block
	prevBlock *Block
	ok        bool
}

var isValidBlockTestCases = []isValidBlockTestCase{
	isValidBlockTestCase{
		&Block{
			Index:        2,
			PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
			Timestamp:    1494177351,
			Data:         "white noise",
		},
		genesisBlock,
		false,
	},
	isValidBlockTestCase{
		&Block{
			Index:        1,
			PreviousHash: "27aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
			Timestamp:    1494177351,
			Data:         "white noise",
		},
		genesisBlock,
		false,
	},
	isValidBlockTestCase{
		&Block{
			Index:        1,
			PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
			Timestamp:    1494177351,
			Data:         "white noise",
		},
		genesisBlock,
		true,
	},
}

func TestIsValidBlock(t *testing.T) {
	for _, testCase := range isValidBlockTestCases {
		ok, err := isValidBlock(testCase.block, testCase.prevBlock)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if ok != testCase.ok {
			t.Errorf("want %t but %t", testCase.ok, ok)
		}
	}
}
