package main

import (
	"os"

	"github.com/0xkhdr/pathframe/internal/adapters/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
