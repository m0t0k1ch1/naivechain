package main

import "testing"

func TestBlockHash(t *testing.T) {
	if genesisBlock.hash() != "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30" {
		t.Errorf("want %q but %q", "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30", genesisBlock.hash())
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
			PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
			Timestamp:    1494177351,
			Data:         "white noise",
		},
		genesisBlock,
		false,
	},
	isValidBlockTestCase{
		&Block{
			Index:        1,
			PreviousHash: "8ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
			Timestamp:    1494177351,
			Data:         "white noise",
		},
		genesisBlock,
		false,
	},
	isValidBlockTestCase{
		&Block{
			Index:        1,
			PreviousHash: "7ca4c614ada5dc59875e7127bbf56083fc4d9ec73f039d3454b09f8891674c30",
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
