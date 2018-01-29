package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

var BlockChain []Block

func (block *Block) generateHash() (string, error) {
	if block.Index <= 0 || block.Timestamp == "" || block.BPM <= 0 {
		return "", fmt.Errorf("Invalid data provided to block")
	}
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	block.Hash = hex.EncodeToString(hashed)
	return hex.EncodeToString(hashed), nil
}

func (block *Block) isBlockValid(oldBlock Block) bool {
	if oldBlock.Index+1 != block.Index {
		return false
	}

	if oldBlock.Hash != block.Hash {
		return false
	}

	expectedHash, err := block.generateHash()
	if expectedHash != block.Hash || err != nil {
		return false
	}
	return true
}

func generateBlock(oldBlock Block, BPM int) (newBlock *Block, err error) {
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash

	newBlock.Hash, err = newBlock.generateHash()

	if err != nil {
		return nil, err
	}

	return newBlock, nil
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(BlockChain) {
		BlockChain = newBlocks
	}
}

func main() {
	fmt.Println("vim-go")

	test := Block{
		Index:     1,
		Timestamp: "gerere",
		BPM:       100,
		Hash:      "",
		PrevHash:  "",
	}

	spew.Dump(test)

	var err error
	test.Hash, err = test.generateHash()

	if err != nil {
		panic(err)
	}

	spew.Dump(test)
}
