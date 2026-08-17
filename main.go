package main

import (
	"os"

	"github.com/anishnarang9/go-chain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
