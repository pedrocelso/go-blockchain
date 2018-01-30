package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Block defines a blockchain block.
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// GenerateHash generates a hash based on the block info and current time
func (block *Block) GenerateHash() (string, error) {
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

// IsBlockValid checks if the block is valid
func (block *Block) IsBlockValid(oldBlock Block) bool {
	if oldBlock.Index+1 != block.Index {
		return false
	}

	if oldBlock.Hash != block.PrevHash {
		return false
	}

	expectedHash, err := block.GenerateHash()
	if expectedHash != block.Hash || err != nil {
		return false
	}
	return true
}
