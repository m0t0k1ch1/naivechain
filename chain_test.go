package main

import "testing"

var (
	genesisBlockHash = "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f"
)

func TestMineBlock(t *testing.T) {
	bc := newBlockchain()

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
	blockHash, err := block.hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != blockHash {
		t.Errorf("want %q but %q", blockHash, latestBlockHash)
	}
}

func TestReplaceChain(t *testing.T) {
	bc := newBlockchain()

	bcNew := newBlockchain()
	if err := bcNew.addBlock(&Block{
		Index:        1,
		PreviousHash: "17aacbe244debc3869a4f604c8136da450283cba3e0467681f398af16871cc3f",
		Timestamp:    1494093545,
		Data:         "white noise",
	}); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	ok, err := bc.tryReplaceBlocks(bcNew)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if !ok {
		t.Fatalf("want %t but %t", true, ok)
	}
	if bc.len() != 2 {
		t.Fatalf("want %d but %d", 2, bc.len())
	}

	latestBlockHash, err := bc.getLatestBlock().hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	blockHash, err := bcNew.getBlock(1).hash()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if latestBlockHash != blockHash {
		t.Errorf("want %q but %q", blockHash, latestBlockHash)
	}
}
