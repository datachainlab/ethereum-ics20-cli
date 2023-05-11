package main

import (
	"log"

	"github.com/datachainlab/ethereum-ics20-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
