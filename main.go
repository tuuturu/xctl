package main

import (
	"os"

	"github.com/deifyed/xctl/cmd"
)

func main() {
	cmd.Execute(os.Stderr)
}
