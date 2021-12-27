package main

import (
	"fmt"
	"os"

	"github.com/ice1n36/brain/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
