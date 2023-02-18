package main

import (
	"fmt"
	"os"
	"tugboat/internal/cli/commands"
)

func main() {

	cli := commands.NewCli()

	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
