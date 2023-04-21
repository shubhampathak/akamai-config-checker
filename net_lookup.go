package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	stagingEdgeKeyStr    = "edgekey-staging"
	productionEdgeKeyStr = "edgekey"
)

var (
	environment string
)

func netLookup(domain string) []net.IP {
	domainEdgeIPs, err := net.LookupIP(domain)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s\n", color.RedString("[Error]"), "Invalid domain provided. Make sure you do not provide an edgeKey domain. Otherwise the lookup may fail.")
		os.Exit(1)
	}
	return domainEdgeIPs
}

func stagLookup(args args, domainCNAME string) []net.IP {
	if strings.Contains(domainCNAME, productionEdgeKeyStr) {
		domainStagingCNAME := strings.Replace(domainCNAME, productionEdgeKeyStr, stagingEdgeKeyStr, 1)
		return netLookup(domainStagingCNAME)
	} else {
		return netLookup(args.domain + "." + stagingEdgeKeyStr + ".net")
	}
}

func prodLookup(args args, domainCNAME string) []net.IP {
	if strings.Contains(domainCNAME, productionEdgeKeyStr) {
		return netLookup(domainCNAME)
	} else {
		return netLookup(args.domain + "." + productionEdgeKeyStr + ".net")
	}
}

func lookup(args args) []net.IP {
	var edgeIPs []net.IP
	domainCNAME, err := net.LookupCNAME(args.domain)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s %v\n", color.RedString("[Error]"), "Invalid domain provided: ", err)
		os.Exit(1)
	}

	switch environment {
	case "staging":
		edgeIPs = stagLookup(args, domainCNAME)
	case "production":
		edgeIPs = prodLookup(args, domainCNAME)
	default:
		edgeIPs = stagLookup(args, domainCNAME)
	}

	return edgeIPs
}
