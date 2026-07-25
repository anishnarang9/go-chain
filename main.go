package main

import (
	"fmt"
	"strconv"

	"github.com/anishnarang9/go-chain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("Second block technically")
	chain.AddBlock("third Block Technically")
	chain.AddBlock(" vava to cinq on cinq back to the vava")
	chain.AddBlock(" Modern Computers are so much faster in like 8 years compared to the tutorial video")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}

}

/// chainhalt
