package main

import "time"

func generateBlock(data string) (*Block, error) {
	previousBlock := blockchain.getLatestBlock()
	previousBlockHash, err := previousBlock.hash()
	if err != nil {
		return nil, err
	}

	return &Block{
		Index:        previousBlock.Index + 1,
		PreviousHash: previousBlockHash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}, nil
}
