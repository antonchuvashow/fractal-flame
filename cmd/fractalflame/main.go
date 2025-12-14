package main

import (
	"os"

	"hw4-fractal-flame/cmd/cli"
)

func main() {
	cmd, err := cli.GetRootCommand()
	if err != nil {
		os.Exit(Broken)
	}

	err = cmd.Execute()
	if err != nil {
		os.Exit(InvalidUsage)
	}
}
