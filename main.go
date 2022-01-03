package main

import (
	"os"

	"github.com/distbuild/worker/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
