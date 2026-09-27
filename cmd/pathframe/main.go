package main

import (
	"context"
	"fmt"
	"os"

	"github.com/0xkhdr/pathframe/internal/adapters/cli"
	mcpadapter "github.com/0xkhdr/pathframe/internal/adapters/mcp"
	"github.com/0xkhdr/pathframe/internal/app"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		if err := mcpadapter.Run(context.Background(), app.Service{Dir: "."}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
