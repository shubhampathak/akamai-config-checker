package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"time"

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

// readFile from the filepath and return file data in []byte and []string
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

// backupHostsFile makes a copy of the existing hosts file before making the changes
func backupHostsFile(args args, orignalHostsFileData []byte, hostsFile string) {
	if args.backup {
		backupFile := hostsFile + time.Now().Format("2006-01-02--15-04-05.bak")
		err := os.WriteFile(backupFile, orignalHostsFileData, 0644)
		if err != nil {
			fmt.Fprintf(color.Output, "%s %s\n %v", color.RedString("[Error]"), "Error creating a backup file.", err)
			os.Exit(1)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("[Success]"), "Backup of the existing hosts file is created as "+backupFile)
		}
	} else if _, err := os.Stat(hostsFile + ".bak"); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(color.Output, "%s %s\n", color.YellowString("[Alert]"), "Creating a backup of the original Hosts file before making changes...")
		err = os.WriteFile(hostsFile+".bak", orignalHostsFileData, 0644)
		if err != nil {
			fmt.Fprintf(color.Output, "%s %s\n %v", color.RedString("[Error]"), "Error creating a backup file.", err)
			os.Exit(1)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("[Success]"), "Backup of the original hosts file is created as "+hostsFile+".bak")
		}
	}
}

// addHostsEntry takes arguments and hostsfile as input and writes entries back to it
func addHostsEntry(args args, hostsFile string) {
	fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Info]"), "Adding hosts entry...")
	orignalHostsFileData, lines := readFile(hostsFile)     // Read existing hosts file
	backupHostsFile(args, orignalHostsFileData, hostsFile) // Backup original hosts file
	updatedLines := make([]string, 0, len(lines))
	var newData string
	if args.domain != "" {
		edgeIPs := lookup(args) // Get Edge Server IPs
		for _, line := range lines {
			if !strings.Contains(line, args.domain) && line != "" {
				updatedLines = append(updatedLines, line)
			}
		}
		updatedLines = append(updatedLines, edgeIPs[0].String()+" "+args.domain)
		newData = strings.Join(updatedLines, "\n") // exclude delimiter for the last line
	} else if args.list != "" {
		_, listLines := readFile(args.list) // Read the list file
		for _, listLine := range listLines {
			if !contains(lines, listLine) && validDomain(listLine) {
				args.domain = listLine
				edgeIPs := lookup(args) // Get Edge Server IPs
				listLine := edgeIPs[0].String() + " " + listLine
				lines = append(lines, listLine)
			}
		}
		// Join the lines back as a string
		newData = strings.Join(lines, "\n") // exclude delimiter for the last line
	} else {
		fmt.Fprintf(color.Output, "%s %s\n", color.BlueString("[Error]"), "Something went wrong! Please try again")
		os.Exit(1)
	}
	// Overwrite the file with the new content
	err := os.WriteFile(linuxHostsFile, []byte(newData), 0644)
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s\n %v", color.RedString("[Error]"), "Error writing file.", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("[Success]"), "Hosts entry was created/updated.")
	}
}
