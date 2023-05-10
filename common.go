package main

import "strings"

// existsStr is a commmon function to check if a target string is found as suffix
// in each element of the array and the elements are not commented with "#"
func existsStr(arrStr []string, targetStr string) bool {
	for _, str := range arrStr {
		if strings.HasSuffix(str, " "+targetStr) && !strings.HasPrefix(str, "#") {
			return true
		}
	}
	return false
}
