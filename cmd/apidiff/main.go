package main

import (
	"fmt"
	"os"

	"github.com/awesomeProject/apidiff/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		if exitErr, ok := err.(*cli.ExitError); ok {
			fmt.Fprintln(os.Stderr, exitErr.Err)
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
