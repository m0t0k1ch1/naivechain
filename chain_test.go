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
	expectedBlockHash, err := block.hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != expectedBlockHash {
		t.Errorf("want %q but %q", expectedBlockHash, latestBlockHash)
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
	expectedBlockHash, err := bcNew.getLatestBlock().hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != expectedBlockHash {
		t.Errorf("want %q but %q", expectedBlockHash, latestBlockHash)
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
				PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
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
				PreviousHash: "27aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
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
				PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
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
