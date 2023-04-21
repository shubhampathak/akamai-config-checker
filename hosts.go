package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

const (
	linuxHostsFile   = "/etc/hosts"
	macOSHostsFile   = "/etc/hosts"
	windowsHostsFile = "C:/Windows/System32/drivers/etc/hosts"
)

// detectOS to find out the operating system being used and return the path to hosts file
func detectOS() (string, error) {
	var hostsFile string
	switch os := runtime.GOOS; os {
	case "windows":
		fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Info]"), "Windows OS detected.")
		hostsFile = windowsHostsFile
	case "linux":
		fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Info]"), "Linux based OS detected.")
		hostsFile = linuxHostsFile
	case "darwin":
		fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Info]"), "macOS detected.")
		hostsFile = macOSHostsFile
	default:
		fmt.Fprintf(color.Output, "%s %s\n", color.RedString("[Error]"), "This OS is not compatible with the tool!")
	}
	return hostsFile, nil
}

func readFile(filePath string) ([]byte, []string) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s %v\n", color.RedString("[Error]"), "Error reading file:", err)
		os.Exit(1)
	}
	// Split the file data line by line
	lines := strings.Split(string(fileData), "\n")
	return fileData, lines
}

func backupHostsFile(orignalHostsFileData []byte, hostsFile string) {
	// Create a backup file if it does not already exist
	if _, err := os.Stat(hostsFile + ".bak"); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(color.Output, "%s %s\n", color.YellowString("[Alert]"), "Creating a backup of the original Hosts file before making changes...")
		err = os.WriteFile(hostsFile+".bak", orignalHostsFileData, 0644)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("[Success]"), "Backup of the original hosts file is created as "+hostsFile+".bak")
		}
	}
}

func addHostsEntry(args args, hostsFile string) {
	edgeIPs := lookup(args) // Get Edge Server IPs
	fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Info]"), "Adding hosts entry...")

	orignalHostsFileData, lines := readFile(hostsFile) // Read existing hosts file
	backupHostsFile(orignalHostsFileData, hostsFile)   // Backup original hosts file
	listLines := readList(args)                        // Read the list file

	// Filter out the lines containing already added domain entry
	filteredLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, args.domain) && line != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	// Append hosts entry <IP> <domain> at the end
	filteredLines = append(filteredLines, edgeIPs[0].String()+" "+args.domain+"\n")

	if len(listLines) > 0 {
		filteredLines = append(filteredLines, listLines...)
	}

	// Join the lines back as a string
	newData := strings.Join(filteredLines, "\n") // exclude delimiter for the last line

	// Overwrite the file with the new content
	err := os.WriteFile(linuxHostsFile, []byte(newData), 0644)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s\n", color.RedString("[Error]"), "Error writing file.")
		os.Exit(1)
	} else {
		fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("[Success]"), "Hosts entry was created.")
	}
}

func readList(args args) []string {
	if args.list != "" {
		edgeIPs := lookup(args)
		_, lines := readFile(args.list)
		updatedLines := make([]string, 0, len(lines))
		for _, line := range lines {
			if line != "" {
				line = edgeIPs[0].String() + " " + line
				updatedLines = append(updatedLines, line)
			}
		}
		return updatedLines
	}
	return nil
}
