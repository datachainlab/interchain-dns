package main

import (
	"os"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/simapp/simd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
