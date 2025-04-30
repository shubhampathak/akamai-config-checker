package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type args struct {
	domain      string // target domain name
	list        string // filepath
	environment string // environment staging or production
	backup      bool   // backup hosts file
}

func parseFlags() args {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d domain] [-l list] [-e environment] [-b]\n", os.Args[0])
		flag.PrintDefaults()
	}

	domain := flag.String("d", "", "Target domain")
	list := flag.String("l", "", "Full Path to list of sub-domains that have same origin IP address as of target domain")
	environment := flag.String("e", "", "Environment against which you want to test. Available options: staging or production")
	backup := flag.Bool("b", false, "Backup the existing hosts file before making changes to it")

	flag.Parse()

	// Ensure at least one of domain or list is provided
	if len(*domain) == 0 && len(*list) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Validate the domain if provided
	if len(*domain) > 0 {
		_, err := isValidDomainAndBehindAkamai(*domain)
		if err != nil {
			fmt.Fprintf(color.Output, "%s %s %v\n", color.RedString("[Error]"), "Domain validation failed:", err)
			os.Exit(1)
		}
		fmt.Fprintf(color.Output, "%s\n", color.BlueString("Target Domain: ")+*domain)
	}

	// Validate the environment
	selectedEnvironment := "staging"
	if len(*environment) > 0 {
		if *environment == "staging" || *environment == "production" {
			selectedEnvironment = *environment
			fmt.Fprintf(color.Output, "%s\n", color.BlueString("Target Environment: ")+selectedEnvironment)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.RedString("[Error]"), "Invalid environment! Available options are: staging or production.")
			os.Exit(1)
		}
	}

	return args{
		domain:      *domain,
		list:        *list,
		environment: selectedEnvironment,
		backup:      *backup,
	}
}
