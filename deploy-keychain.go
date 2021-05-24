package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nint8835/deploy-keychain/deploykeychain"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("This tool is not intended to be ran directly.")
		fmt.Println("See https://github.com/nint8835/deploy-keychain/README.md for usage and configuration details.")
		return
	}

	config, err := deploykeychain.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
	}

	deploykeychain.Log(fmt.Sprintf("Args: %s", os.Args))

	account, repository, err := deploykeychain.IdentifyRepository(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to determine what repository is being used.")
		os.Exit(1)
	}

	keyFile, err := deploykeychain.DetermineKeyFile(config, account, repository)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine key: %s\n", err)
		os.Exit(1)
	}

	args := append([]string{"-i", keyFile}, os.Args[1:]...)
	deploykeychain.Log(fmt.Sprintf("Running %s with args %+v", config.SSHCommand, args))

	cmd := exec.Command(config.SSHCommand, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running SSH: %s\n", err)
		os.Exit(1)
	}
}
