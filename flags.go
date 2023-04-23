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
	environment string // environment
}

func parseFlags() args {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d domain] [-l list] [-e environment]\n", os.Args[0])
		flag.PrintDefaults()
	}

	domain := flag.String("d", "", "Target domain")
	list := flag.String("l", "", "Full Path to list of sub-domains that have same origin IP address as of target domain")
	environment := flag.String("e", "", "Environment against which you want to test. Available options: staging or production")

	flag.Parse()

	if len(*domain) == 0 && len(*list) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(*domain) > 0 {
		if validDomain(*domain) {
			fmt.Fprintf(color.Output, "%s\n", color.BlueString("Target Domain: ")+*domain)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.RedString("[Error]"), "Invalid input! Please provide a valid domain.")
			os.Exit(1)
		}
	}

	var allowedEnvironment string
	if *environment == "staging" || *environment == "production" {
		allowedEnvironment = *environment
		fmt.Fprintf(color.Output, "%s\n", color.BlueString("Target Environment: ")+allowedEnvironment)
	}

	return args{
		domain:      *domain,
		list:        *list,
		environment: allowedEnvironment,
	}
}
