package main

import (
	"crypto/sha256"
	"fmt"
)

var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
	Timestamp:    1465154705,
	Data:         "my genesis block!!",
}

type Block struct {
	Index        int64  `json:"index"`
	PreviousHash string `json:"previousHash"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
}

func (block *Block) hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf(
		"%d%s%d%s",
		block.Index, block.PreviousHash, block.Timestamp, block.Data,
	))))
}

func isValidBlock(block, prevBlock *Block) bool {
	return block.Index == prevBlock.Index+1 &&
		block.PreviousHash == prevBlock.hash()
}
