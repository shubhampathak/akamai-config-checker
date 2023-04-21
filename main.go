package main

import (
	"log"
)

func main() {
	args := parseFlags()
	hostsFile, err := detectOS()
	if err != nil || hostsFile == "" {
		log.Fatalln(err)
	} else {
		addHostsEntry(args, hostsFile)
	}
}
