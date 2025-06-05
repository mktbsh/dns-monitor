package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Domains      []string
	RecordType   string
	Interval     time.Duration
	Servers      []string
	AllServers   bool
	UntilChange  bool
	OutputFile   string
	NoColor      bool
	ShowHelp     bool
	ShowVersion  bool
}

func ParseArgs(args []string) (*Config, error) {
	config := &Config{
		RecordType: "A",
		Interval:   5 * time.Second,
		Servers:    []string{},
	}

	i := 1
	for i < len(args) {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			config.ShowHelp = true
			return config, nil
		case arg == "-v" || arg == "--version":
			config.ShowVersion = true
			return config, nil
		case arg == "-t" || arg == "--type":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("option %s requires a value", arg)
			}
			config.RecordType = strings.ToUpper(args[i+1])
			i += 2
		case arg == "-i" || arg == "--interval":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("option %s requires a value", arg)
			}
			duration, err := parseDuration(args[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid interval: %v", err)
			}
			config.Interval = duration
			i += 2
		case arg == "-s" || arg == "--server":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("option %s requires a value", arg)
			}
			server := args[i+1]
			if !strings.Contains(server, ":") {
				server += ":53"
			}
			config.Servers = append(config.Servers, server)
			i += 2
		case arg == "--all-servers":
			config.AllServers = true
			i++
		case arg == "--until-change":
			config.UntilChange = true
			i++
		case arg == "-o" || arg == "--output":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("option %s requires a value", arg)
			}
			config.OutputFile = args[i+1]
			i += 2
		case arg == "--no-color":
			config.NoColor = true
			i++
		case strings.HasPrefix(arg, "-"):
			return nil, fmt.Errorf("unknown option: %s", arg)
		default:
			config.Domains = append(config.Domains, args[i:]...)
			i = len(args)
		}
	}

	if config.AllServers {
		config.Servers = []string{"8.8.8.8:53", "1.1.1.1:53", "1.0.0.1:53"}
	}

	if len(config.Domains) == 0 && !config.ShowHelp && !config.ShowVersion {
		return nil, fmt.Errorf("at least one domain must be specified")
	}

	if !isValidRecordType(config.RecordType) {
		return nil, fmt.Errorf("unsupported record type: %s", config.RecordType)
	}

	return config, nil
}

func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "s") {
		seconds, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(seconds) * time.Second, nil
	}
	if strings.HasSuffix(s, "m") {
		minutes, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(minutes) * time.Minute, nil
	}
	if strings.HasSuffix(s, "h") {
		hours, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(hours) * time.Hour, nil
	}

	seconds, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s (use formats like 5s, 2m, 1h)", s)
	}
	return time.Duration(seconds) * time.Second, nil
}

func isValidRecordType(recordType string) bool {
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT"}
	for _, t := range validTypes {
		if t == recordType {
			return true
		}
	}
	return false
}

func (c *Config) Print() {
	fmt.Printf("Domains: %v\n", c.Domains)
	fmt.Printf("Record Type: %s\n", c.RecordType)
	fmt.Printf("Interval: %s\n", c.Interval)
	fmt.Printf("Servers: %v\n", c.Servers)
	fmt.Printf("Until Change: %t\n", c.UntilChange)
	if c.OutputFile != "" {
		fmt.Printf("Output File: %s\n", c.OutputFile)
	}
	fmt.Printf("No Color: %t\n", c.NoColor)
}