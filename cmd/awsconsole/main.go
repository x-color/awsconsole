package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/x-color/awsconsole"
)

func main() {
	update := flag.Bool("update", false, "update accounts and roles file")
	flag.Parse()

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	var err error
	if *update {
		err = awsconsole.GenerateAccountsRolesFile(profile)
	} else {
		err = awsconsole.Jump(profile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
