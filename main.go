package main

import (
	"os"

	"github.com/Ashtacore/golang_blockchain_demo/cli"
)

func main() {
	defer os.Exit(0)
	CLI := cli.CommandLine{}
	CLI.Run()
}
