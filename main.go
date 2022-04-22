package main

import (
	"fmt"
	"strconv"

	"github.com/Ashtacore/golang_blockchain_demo/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("Block Two")
	chain.AddBlock("Block Three")
	chain.AddBlock("Block Four")

	for _, block := range chain.Blocks {
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
