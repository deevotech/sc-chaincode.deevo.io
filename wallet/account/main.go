package main

import (
	"os"
	//"github.com/spf13/cobra"
)

func main() {
	if err := RunMain(os.Args); err != nil {
		os.Exit(1)
	}
}

// RunMain is the fabric-ca server main
func RunMain(args []string) error {
	// Save the os.Args
	saveOsArgs := os.Args
	os.Args = args

	cmdName := ""
	if len(args) > 1 {
		cmdName = args[1]
	}
	scmd := NewCommand(cmdName)

	// Execute the command
	err := scmd.Execute()

	// Restore original os.Args
	os.Args = saveOsArgs

	return err
}
