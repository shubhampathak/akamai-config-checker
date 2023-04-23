# Akamai Configuration Checker
A command line tool to validate configuration on Akamai staging or production networks. This tool automates the process of finding staging/production Akamai Edge Hostname IPs for a provided domain(s) and adds the entries to the system's host file.

## Supported OS
This tool should work on:

- Linux based
- Windows 
- macOS

## Features

This tool will:

1. Find out the staging or production Akamai Edge Hostname of a given domain or list of domains.
2. After that, it will look for the IPs assigned to the discovered EdgeKey Hostname.
[Note: it only supports ".edgekey.net" based hostnames for now.]
3. If a subdomain is using root domain's edgekey hostname, then it will query for the root domain's edge hostname instead.
4. Upon finding the IP(s), it will first backup the original hosts file and add the entry[ies] to the system's hosts file.

## Installation

Download the pre-built binary from the [releases](https://github.com/shubhampathak/akamai-config-checker/releases) page.

or 

Use the below command to directly install it if you have the recent compiler:

```bash
go install -v github.com/shubhampathak/akamai-config-checker@latest
```

or 

Manually build it if Go is already installed in your system:
``` bash
git clone https://github.com/shubhampathak/akamai-config-checker.git
cd akamai-config-checker
go build
```

## Usage

**Note: Since it makes changes to the system's hosts file, hence running the tool as sudo/admin user privilege is required.**

1. Linux/macOS users: 
```bash
sudo ./akamai-config-checker -h
```

2. Windows users:
Open Powershell or command prompt as administrator and then run the executable.
```
cd <path_to_downloaded_binary>
akamai-config-checker -h
```


```
Usage: sudo ./akamai-config-checker [-d domain] [-l list] [-e environment]
  -d string
        Target domain
  -e string
        Environment against which you want to test. Available options: staging or production
  -l string
        Full Path to list of sub-domains that have same origin IP address as of target domain
```

1. To add entry for a Single Domain:
``` bash
sudo ./akamai-config-checker -d example.com
```

2. To add entry for a list of domains/sub-domains:
``` bash
sudo ./akamai-config-checker -l <path_to_list_file>
```

3. Use `-e` environment flag `-e staging` if you want to test for staging or `-e production` if you want to test on production environment. If `-e` flag is not provided, then it will check for staging by default.
``` bash
sudo ./akamai-config-checker -d example.com -e staging
```
``` bash
sudo ./akamai-config-checker -l <path_to_list_file> -e production
```

Please note that if both `-d` (domain) and `-l` (list) flags are provided, then it will only use the `-d` flag. The list flag `-l` will be ignored.

## License

[MIT](https://github.com/shubhampathak/akamai-config-checker/blob/main/LICENSE)
## Disclaimer:

This tool makes changes to system's hosts file, thus requires root privilege to run. It makes a backup of the original hosts file in the same directory. This repository's maintainer shall not be liable for any damage caused to the system or data due to the usage of this tool. 

