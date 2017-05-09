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

	block := bc.generateBlock("white noise")
	if block.Index != bc.getLatestBlock().Index+1 {
		t.Errorf("want %d but %d", bc.getLatestBlock().Index+1, block.Index)
	}
	if block.Data != "white noise" {
		t.Errorf("want %q but %q", "white noise", block.Data)
	}
	if block.PreviousHash != bc.getLatestBlock().hash() {
		t.Errorf("want %q but %q", bc.getLatestBlock().hash(), block.PreviousHash)
	}
}

func TestAddBlock(t *testing.T) {
	bc := newTestBlockchain([]*Block{genesisBlock})
	block := &Block{
		Index:        1,
		PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
		Timestamp:    1494093545,
		Data:         "white noise",
	}

	bc.addBlock(block)
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}
	if bc.getLatestBlock().hash() != block.hash() {
		t.Errorf("want %q but %q", block.hash(), bc.getLatestBlock().hash())
	}
}

func TestReplaceBlocks(t *testing.T) {
	bc := newTestBlockchain([]*Block{genesisBlock})
	bcNew := newTestBlockchain([]*Block{
		genesisBlock,
		&Block{
			Index:        1,
			PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
			Timestamp:    1494093545,
			Data:         "white noise",
		},
	})

	bc.replaceBlocks(bcNew)
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}
	if bc.getLatestBlock().hash() != bcNew.getLatestBlock().hash() {
		t.Errorf("want %q but %q", bcNew.getLatestBlock().hash(), bc.getLatestBlock().hash())
	}
}

type isValidChainTestCase struct {
	blockchain *Blockchain
	ok         bool
}

var isValidChainTestCases = []isValidChainTestCase{
	isValidChainTestCase{
		newTestBlockchain([]*Block{}),
		false,
	},
	isValidChainTestCase{
		newTestBlockchain([]*Block{
			&Block{
				Index:        0,
				PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
				Timestamp:    1465154705,
				Data:         "bad genesis block!!",
			},
		}),
		false,
	},
	isValidChainTestCase{
		newTestBlockchain([]*Block{
			genesisBlock,
			&Block{
				Index:        2,
				PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
				Timestamp:    1494177351,
				Data:         "white noise",
			},
		}),
		false,
	},
	isValidChainTestCase{
		newTestBlockchain([]*Block{
			genesisBlock,
			&Block{
				Index:        1,
				PreviousHash: "8ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
				Timestamp:    1494177351,
				Data:         "white noise",
			},
		}),
		false,
	},
	isValidChainTestCase{
		newTestBlockchain([]*Block{
			genesisBlock,
			&Block{
				Index:        1,
				PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
				Timestamp:    1494177351,
				Data:         "white noise",
			},
		}),
		true,
	},
}

func TestIsValidChain(t *testing.T) {
	for _, testCase := range isValidChainTestCases {
		ok, err := testCase.blockchain.isValid()
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if ok != testCase.ok {
			t.Errorf("want %t but %t", testCase.ok, ok)
		}
	}
}
