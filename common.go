package main

import "strings"

// contains is a commmon function to check if an input string is found in an array of string
func contains(as []string, s string) bool {
	for _, a := range as {
		if strings.Contains(a, s) {
			return true
		}
	}
	return false
}
