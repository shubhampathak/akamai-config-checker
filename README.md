# Akamai Configuration Checker
A command-line tool that helps you validate configuration on Akamai staging and production networks. This tool automates the procedure for locating and adding Akamai Edge Hostname IPs for a specified domain(s) and then adding them to the system's host file. In short, it automates steps 2-4 mentioned in the respective [Akamai help document](https://techdocs.akamai.com/api-acceleration/docs/test-stage).

![akamai-config-checker](https://github.com/shubhampathak/akamai-config-checker/assets/20816337/f6b29f70-8895-4a66-8f24-c350ab37a6ba)
## Supported OS
This tool should work on the below-mentioned Operating Systems:

- Linux based
- Windows 
- macOS

## Features

This tool will perform the following tasks:

1. Find out the staging or production Akamai Edge Hostname of a given domain or a list of domain names.
2. After that, it will look for the IPs assigned to the discovered EdgeKey Hostname. [Note: it only supports ".edgekey.net" based hostnames for now.]
3. If a subdomain utilises the root domain's edge key hostname, it will query for the root domain's edge hostname instead.
4. Upon finding the IP(s), it will back up the original hosts file and add the entry[ies] to the system's hosts file.

## Installation

[Recommended] Use the below command to directly install it if you have the recent compiler:

```bash
go install -v github.com/shubhampathak/akamai-config-checker@latest
```

or

Download the pre-built binary from the [releases](https://github.com/shubhampathak/akamai-config-checker/releases) page. 

or 

Manually build it if Go is already installed in your system:
``` bash
git clone https://github.com/shubhampathak/akamai-config-checker.git
cd akamai-config-checker
go build
```

## Usage

**Note: Since it modifies the system's hosts file, that's why, running the tool as sudo/admin user privilege is required.**

1. Linux/macOS users: 
```bash
sudo ./akamai-config-checker -h
```

2. Windows users:
Open Powershell or command prompt as administrator and then run the executable.
```
cd <path_to_downloaded_binary>
.\akamai-config-checker.exe -h
```

```
Usage: ./akamai-config-checker [-d domain] [-l list] [-e environment] [-b]
  -b    Backup the existing hosts file before making changes to it
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

A list file can be a file that contains one domain/sub-domain per line:
```
example.com
sub.example.com
multi.level.example.edu
test-api.example.org
```

3. Use the -b (backup) flag to back up the existing hosts file. If this flag is absent, the tool will automatically make a backup of the original hosts file as hosts.bak only if the hosts.bak is not already present in the hosts file directory. The backup flag -b always generates a new backup file hosts--<date_with_time>.bak from the existing hosts file.

``` bash
sudo ./akamai-config-checker -l <path_to_list_file> -b
```

4. Use the -e (environment) flag `-e staging` for testing on Akamai staging or `-e production` for testing on the Akamai production environment. If `-e` is not supplied, it will check for the staging environment as default.
``` bash
sudo ./akamai-config-checker -d example.com -e staging
```
``` bash
sudo ./akamai-config-checker -l <path_to_list_file> -e production
```

Please note that if both -d (domain) and -l (list) flags are present, then the input of the -l (list) flag will be ignored.

## License

[MIT](https://github.com/shubhampathak/akamai-config-checker/blob/main/LICENSE)
## Disclaimer:

This tool modifies the system's hosts file, thus requiring root privilege to run. It makes a backup of the original hosts file in the same directory. This repository's maintainer shall not be liable for any damage caused to the system or data due to the usage of this tool.

