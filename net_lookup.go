package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	stagingEdgeKeySuffix    = ".edgekey-staging.net"
	productionEdgeKeySuffix = ".edgekey.net"
)

func validDomain(domain string) bool {
	hostNames, err := net.LookupHost(domain)
	if err != nil || len(hostNames) == 0 {
		return false
	}
	return true
}

func netLookup(domain string) ([]net.IP, error) {
	domainEdgeIPs, err := net.LookupIP(domain)
	return domainEdgeIPs, err
}

func stagLookup(domain, domainCNAME string) ([]net.IP, error) {
	if strings.Contains(domainCNAME, productionEdgeKeySuffix) {
		domainStagingCNAME := strings.Replace(domainCNAME, productionEdgeKeySuffix, stagingEdgeKeySuffix, 1)
		return netLookup(domainStagingCNAME)
	} else {
		domainEdgeIPs, err := netLookup(domain + stagingEdgeKeySuffix)
		if err != nil {
			parts := strings.Split(domain, ".")
			rootDomain := strings.Join(parts[len(parts)-2:], ".")
			domainEdgeIPs, err := net.LookupIP(rootDomain + stagingEdgeKeySuffix)
			if err != nil {
				return nil, err
			}
			return domainEdgeIPs, nil
		}
		return domainEdgeIPs, nil
	}
}

func prodLookup(domain, domainCNAME string) ([]net.IP, error) {
	if strings.Contains(domainCNAME, productionEdgeKeySuffix) {
		return netLookup(domainCNAME)
	} else {
		domainEdgeIPs, err := netLookup(domain + productionEdgeKeySuffix)
		if err != nil {
			parts := strings.Split(domain, ".")
			rootDomain := strings.Join(parts[len(parts)-2:], ".")
			domainEdgeIPs, err := net.LookupIP(rootDomain + productionEdgeKeySuffix)
			if err != nil {
				return nil, err
			}
			return domainEdgeIPs, nil
		}
		return domainEdgeIPs, nil
	}
}

func lookup(args args) []net.IP {
	var edgeIPs []net.IP
	domainCNAME, err := net.LookupCNAME(args.domain)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s %v\n", color.RedString("[Error]"), "Invalid domain provided! Malformed input domain or list provided. ", err)
		os.Exit(1)
	}
	switch args.environment {
	case "staging":
		edgeIPs, _ = stagLookup(args.domain, domainCNAME)
	case "production":
		edgeIPs, _ = prodLookup(args.domain, domainCNAME)
	default:
		edgeIPs, _ = stagLookup(args.domain, domainCNAME)
	}
	return edgeIPs
}
