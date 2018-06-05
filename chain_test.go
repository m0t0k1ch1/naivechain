package main

import (
	"sync"
	"testing"
)

func newTestBlockchain(blocks Blocks) *Blockchain {
	return &Blockchain{
		blocks: blocks,
		mu:     sync.RWMutex{},
	}
}

func TestGenerateBlock(t *testing.T) {
	bc := newTestBlockchain(Blocks{testGenesisBlock})

	block := bc.generateBlock("white noise")
	if block.Index != bc.getLatestBlock().Index+1 {
		t.Errorf("want %d but %d", bc.getLatestBlock().Index+1, block.Index)
	}
	if block.Data != "white noise" {
		t.Errorf("want %q but %q", "white noise", block.Data)
	}
	if block.PreviousHash != bc.getLatestBlock().Hash {
		t.Errorf("want %q but %q", bc.getLatestBlock().Hash, block.PreviousHash)
	}
}

func TestAddBlock(t *testing.T) {
	bc := newTestBlockchain(Blocks{testGenesisBlock})
	block := &Block{
		Index:        1,
		PreviousHash: testGenesisBlock.Hash,
		Timestamp:    1494177351,
		Data:         "white noise",
		Hash:         "1cee23ac6ce3589aedbd92213e0dbf8ab41f8f8e6181a92c1a8243df4b32078b",
	}

	bc.addBlock(block)
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}
	if bc.getLatestBlock().Hash != block.Hash {
		t.Errorf("want %q but %q", block.Hash, bc.getLatestBlock().Hash)
	}
}

func TestReplaceBlocks(t *testing.T) {
	bc := newTestBlockchain(Blocks{testGenesisBlock})
	blocks := Blocks{
		testGenesisBlock,
		&Block{
			Index:        1,
			PreviousHash: testGenesisBlock.Hash,
			Timestamp:    1494093545,
			Data:         "white noise",
			Hash:         "1cee23ac6ce3589aedbd92213e0dbf8ab41f8f8e6181a92c1a8243df4b32078b",
		},
	}

	bc.replaceBlocks(blocks)
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}
	if bc.getLatestBlock().Hash != blocks[len(blocks)-1].Hash {
		t.Errorf("want %q but %q", blocks[len(blocks)-1].Hash, bc.getLatestBlock().Hash)
	}
}

func TestIsValidChain(t *testing.T) {
	testCases := []struct {
		name       string
		blockchain *Blockchain
		ok         bool
	}{
		{
			"empty",
			newTestBlockchain(Blocks{}),
			false,
		},
		{
			"invalid genesis block",
			newTestBlockchain(Blocks{
				&Block{
					Index:        0,
					PreviousHash: "0",
					Timestamp:    1465154705,
					Data:         "bad genesis block!!",
					Hash:         "627ab16dbcede0cfa91c85a88c30c4eaae41b8500a961d0d09451323c6e25bf8",
				},
			}),
			false,
		},
		{
			"invalid block",
			newTestBlockchain(Blocks{
				testGenesisBlock,
				&Block{
					Index:        2,
					PreviousHash: testGenesisBlock.Hash,
					Timestamp:    1494177351,
					Data:         "white noise",
					Hash:         "6e27d73b81b2abf47e6766b8aad12a114614fccac669d0d2162cb842f0484420",
				},
			}),
			false,
		},
		{
			"valid",
			newTestBlockchain(Blocks{
				testGenesisBlock,
				&Block{
					Index:        1,
					PreviousHash: testGenesisBlock.Hash,
					Timestamp:    1494177351,
					Data:         "white noise",
					Hash:         "1cee23ac6ce3589aedbd92213e0dbf8ab41f8f8e6181a92c1a8243df4b32078b",
				},
			}),
			true,
		},
	}

	for _, tc := range testCases {
		if ok := tc.blockchain.isValid(); ok != tc.ok {
			t.Errorf("[%s] want %t but %t", tc.name, tc.ok, ok)
		}
	}
}
