package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	prevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.prevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]

}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}
func Genesis() *Block {
	return CreateBlock("genesis", []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{(Genesis())}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("Second block technically")
	chain.AddBlock("third Block Technically")
	chain.AddBlock(" vava to cinq on cinq back to the vava")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.prevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}

}

/// chainhalt
