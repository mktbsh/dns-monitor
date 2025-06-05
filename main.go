package main

import (
	"fmt"
	"log"
	"os"
)

const Version = "1.0.0"

func main() {
	config, err := ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if config.ShowHelp {
		printUsage()
		return
	}

	if config.ShowVersion {
		fmt.Printf("DNS Monitor Tool v%s\n", Version)
		return
	}

	monitor := NewMonitor(config)
	if err := monitor.Start(); err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `DNS Monitor Tool v%s

USAGE:
    dns-monitor [OPTIONS] DOMAIN [DOMAIN...]

OPTIONS:
    -t, --type TYPE          DNS record type (A, AAAA, CNAME, MX, TXT, etc.) [default: A]
    -i, --interval DURATION  Check interval [default: 5s]
    -s, --server SERVER      Specify DNS server (multiple allowed)
    --all-servers           Query all major DNS servers (8.8.8.8, 1.1.1.1, 1.0.0.1)
    --until-change          Monitor until change mode
    -o, --output FILE       Log file output destination
    --no-color              Disable colored output
    -h, --help              Display help
    -v, --version           Display version

EXAMPLES:
    dns-monitor example.com
    dns-monitor -i 30s example.com api.example.com
    dns-monitor -t CNAME --until-change www.example.com
    dns-monitor -s 8.8.8.8 -s 1.1.1.1 example.com
    dns-monitor -o /var/log/dns-monitor.log example.com
`, Version)
}