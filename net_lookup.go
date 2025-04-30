package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	stagingEdgeKeySuffix      = ".edgekey-staging.net"
	productionEdgeKeySuffix   = ".edgekey.net"
	stagingEdgeSuiteSuffix    = ".edgesuite-staging.net"
	productionEdgeSuiteSuffix = ".edgesuite.net"
)

// Perform a DNS lookup for a domain
func netLookup(domain string) ([]net.IP, error) {
	return net.LookupIP(domain)
}

// Helper function to handle staging and production lookups
func edgeLookup(domain, domainCNAME, suffix string) ([]net.IP, error) {
	var targetDomain string
	// Check if the CNAME contains either production suffix
	if strings.Contains(domainCNAME, productionEdgeKeySuffix) {
		targetDomain = strings.Replace(domainCNAME, productionEdgeKeySuffix, suffix, 1)
	} else if strings.Contains(domainCNAME, productionEdgeSuiteSuffix) {
		targetDomain = strings.Replace(domainCNAME, productionEdgeSuiteSuffix, suffix, 1)
	} else {
		targetDomain = domain + suffix
	}
	// Perform the DNS lookup
	ips, err := netLookup(targetDomain)
	if err != nil {
		// Handle cases where the root domain needs to be queried
		parts := strings.Split(domain, ".")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid domain format")
		}
		rootDomain := strings.Join(parts[len(parts)-2:], ".")
		ips, err = netLookup(rootDomain + suffix)
		if err != nil {
			return nil, err
		}
	}
	return validateIPs(ips)
}

// Lookup IPs for staging or production environments
func lookup(args args) ([]net.IP, error) {
	// Determine the suffix based on the environment
	var suffix string
	switch args.environment {
	case "staging":
		// Use staging suffixes
		if strings.Contains(args.domain, ".edgesuite.net") {
			suffix = stagingEdgeSuiteSuffix
		} else {
			suffix = stagingEdgeKeySuffix
		}
	case "production":
		// Use production suffixes
		if strings.Contains(args.domain, ".edgesuite.net") {
			suffix = productionEdgeSuiteSuffix
		} else {
			suffix = productionEdgeKeySuffix
		}
	default:
		return nil, fmt.Errorf("unsupported environment: %s", args.environment)
	}

	// Perform the edge lookup
	domainCNAME, err := net.LookupCNAME(args.domain)
	if err != nil {
		return nil, fmt.Errorf("invalid domain provided: %v", err)
	}
	return edgeLookup(args.domain, domainCNAME, suffix)
}
