package main

import (
	"sort"
	"testing"
)

var testGenesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1465154705,
	Data:         "my genesis block!!",
	Hash:         "816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
}

func TestSortBlocks(t *testing.T) {
	blocks := Blocks{
		&Block{
			Index: 2,
		},
		&Block{
			Index: 3,
		},
		&Block{
			Index: 1,
		},
	}

	sort.Sort(blocks)

	var i int64 = 1
	for _, block := range blocks {
		if block.Index != i {
			t.Errorf("want %d but %d", i, block.Index)
		}
		i++
	}
}

func TestBlockHash(t *testing.T) {
	if testGenesisBlock.hash() != testGenesisBlock.Hash {
		t.Errorf("want %q but %q", testGenesisBlock.Hash, testGenesisBlock.hash())
	}
}

func TestIsValidBlock(t *testing.T) {
	testCases := []struct {
		name      string
		block     *Block
		prevBlock *Block
		ok        bool
	}{
		{
			"invalid index",
			&Block{
				Index:        2,
				PreviousHash: testGenesisBlock.Hash,
				Timestamp:    1494177351,
				Data:         "white noise",
				Hash:         "6e27d73b81b2abf47e6766b8aad12a114614fccac669d0d2162cb842f0484420",
			},
			testGenesisBlock,
			false,
		},
		{
			"invalid previous hash",
			&Block{
				Index:        1,
				PreviousHash: "016534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
				Timestamp:    1494177351,
				Data:         "white noise",
				Hash:         "03bf0215fef25dbf56e7b26ac57f7412cd10aea5e9f2bd8056a349bfaa15bfa5",
			},
			testGenesisBlock,
			false,
		},
		{
			"invalid hash",
			&Block{
				Index:        1,
				PreviousHash: testGenesisBlock.Hash,
				Timestamp:    1494177351,
				Data:         "white noise",
				Hash:         testGenesisBlock.Hash,
			},
			testGenesisBlock,
			false,
		},
		{
			"valid",
			&Block{
				Index:        1,
				PreviousHash: testGenesisBlock.Hash,
				Timestamp:    1494177351,
				Data:         "white noise",
				Hash:         "1cee23ac6ce3589aedbd92213e0dbf8ab41f8f8e6181a92c1a8243df4b32078b",
			},
			testGenesisBlock,
			true,
		},
	}

	for _, tc := range testCases {
		if ok := isValidBlock(tc.block, tc.prevBlock); ok != tc.ok {
			t.Errorf("[%s] want %t but %t", tc.name, tc.ok, ok)
		}
	}
}
