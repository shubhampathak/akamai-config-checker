package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// DNS cache to reduce redundant lookups
var dnsCache = cache.New(5*time.Minute, 10*time.Minute)

// isValidDomainName verifies if a domain name is valid
func isValidDomainName(domain string) bool {
	var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}

// cachedLookupCNAME performs a CNAME lookup with caching
func cachedLookupCNAME(domain string) (string, error) {
	if cname, found := dnsCache.Get(domain); found {
		return cname.(string), nil
	}
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return "", fmt.Errorf("failed to perform CNAME lookup: %v", err)
	}
	dnsCache.Set(domain, cname, cache.DefaultExpiration)
	return cname, nil
}

// isDomainBehindAkamai checks if a domain is behind Akamai by performing a CNAME lookup
func isDomainBehindAkamai(domain string) (bool, error) {
	var cname string
	var err error
	for i := range 3 {
		cname, err = cachedLookupCNAME(domain)
		if err == nil {
			break
		}
		log.Printf("[Retry %d] Failed to lookup CNAME for %s: %v\n", i+1, domain, err)
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		return false, fmt.Errorf("failed to perform CNAME lookup after retries: %v", err)
	}
	akamaiEdgeDomainsPatterns := []string{".edgekey.net", ".edgesuite.net"}
	for _, domainPattern := range akamaiEdgeDomainsPatterns {
		if strings.Contains(cname, domainPattern) {
			return true, nil
		}
	}
	return false, nil
}

// isValidDomainAndBehindAkamai checks if a domain is valid and behind Akamai
// by performing a DNS lookup and validating the format
func isValidDomainAndBehindAkamai(domain string) (bool, error) {
	// Validate the domain format
	if !isValidDomainName(domain) {
		return false, fmt.Errorf("invalid domain format")
	}
	// Check if the domain resolves to at least one hostname
	hostNames, err := net.LookupHost(domain)
	if err != nil || len(hostNames) == 0 {
		return false, fmt.Errorf("domain does not resolve to any hostnames")
	}
	// Check if the domain is behind Akamai
	isAkamai, err := isDomainBehindAkamai(domain)
	if err != nil {
		return false, fmt.Errorf("failed to check if domain is behind Akamai: %v", err)
	}
	if !isAkamai {
		return false, fmt.Errorf("domain is probably not behind Akamai")
	}
	return true, nil
}

// isPublicIP validates if an IP is public
func isPublicIP(ip net.IP) bool {
	return !(ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified())
}

// validateIPs filters and returns only public IPs
func validateIPs(ips []net.IP) ([]net.IP, error) {
	var publicIPs []net.IP
	for _, ip := range ips {
		if isPublicIP(ip) {
			publicIPs = append(publicIPs, ip)
		}
	}
	if len(publicIPs) == 0 {
		return nil, fmt.Errorf("no public IPs found")
	}
	return publicIPs, nil
}
