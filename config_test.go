package main

import (
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    *Config
		expectError bool
	}{
		{
			name: "basic domain only",
			args: []string{"dns-monitor", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{},
			},
		},
		{
			name: "multiple domains",
			args: []string{"dns-monitor", "example.com", "test.com"},
			expected: &Config{
				Domains:    []string{"example.com", "test.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{},
			},
		},
		{
			name: "custom record type",
			args: []string{"dns-monitor", "-t", "CNAME", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "CNAME",
				Interval:   5 * time.Second,
				Servers:    []string{},
			},
		},
		{
			name: "custom interval",
			args: []string{"dns-monitor", "-i", "30s", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   30 * time.Second,
				Servers:    []string{},
			},
		},
		{
			name: "custom server",
			args: []string{"dns-monitor", "-s", "8.8.8.8", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{"8.8.8.8:53"},
			},
		},
		{
			name: "multiple servers",
			args: []string{"dns-monitor", "-s", "8.8.8.8", "-s", "1.1.1.1:53", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{"8.8.8.8:53", "1.1.1.1:53"},
			},
		},
		{
			name: "all servers flag",
			args: []string{"dns-monitor", "--all-servers", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{"8.8.8.8:53", "1.1.1.1:53", "1.0.0.1:53"},
				AllServers: true,
			},
		},
		{
			name: "until change mode",
			args: []string{"dns-monitor", "--until-change", "example.com"},
			expected: &Config{
				Domains:     []string{"example.com"},
				RecordType:  "A",
				Interval:    5 * time.Second,
				Servers:     []string{},
				UntilChange: true,
			},
		},
		{
			name: "output file",
			args: []string{"dns-monitor", "-o", "/tmp/dns.log", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{},
				OutputFile: "/tmp/dns.log",
			},
		},
		{
			name: "no color",
			args: []string{"dns-monitor", "--no-color", "example.com"},
			expected: &Config{
				Domains:    []string{"example.com"},
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{},
				NoColor:    true,
			},
		},
		{
			name: "help flag",
			args: []string{"dns-monitor", "--help"},
			expected: &Config{
				RecordType: "A",
				Interval:   5 * time.Second,
				Servers:    []string{},
				ShowHelp:   true,
			},
		},
		{
			name: "version flag",
			args: []string{"dns-monitor", "-v"},
			expected: &Config{
				RecordType:  "A",
				Interval:    5 * time.Second,
				Servers:     []string{},
				ShowVersion: true,
			},
		},
		{
			name:        "no domains",
			args:        []string{"dns-monitor"},
			expectError: true,
		},
		{
			name:        "invalid record type",
			args:        []string{"dns-monitor", "-t", "INVALID", "example.com"},
			expectError: true,
		},
		{
			name:        "missing interval value",
			args:        []string{"dns-monitor", "-i"},
			expectError: true,
		},
		{
			name:        "invalid interval",
			args:        []string{"dns-monitor", "-i", "invalid", "example.com"},
			expectError: true,
		},
		{
			name:        "unknown option",
			args:        []string{"dns-monitor", "--unknown", "example.com"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArgs(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !configsEqual(config, tt.expected) {
				t.Errorf("config mismatch.\nExpected: %+v\nGot: %+v", tt.expected, config)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{"seconds with suffix", "30s", 30 * time.Second, false},
		{"minutes with suffix", "5m", 5 * time.Minute, false},
		{"hours with suffix", "2h", 2 * time.Hour, false},
		{"seconds without suffix", "45", 45 * time.Second, false},
		{"invalid format", "abc", 0, true},
		{"invalid number", "5x", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDuration(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsValidRecordType(t *testing.T) {
	tests := []struct {
		recordType string
		expected   bool
	}{
		{"A", true},
		{"AAAA", true},
		{"CNAME", true},
		{"MX", true},
		{"TXT", true},
		{"NS", false},
		{"SOA", false},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.recordType, func(t *testing.T) {
			result := isValidRecordType(tt.recordType)
			if result != tt.expected {
				t.Errorf("expected %v for %s, got %v", tt.expected, tt.recordType, result)
			}
		})
	}
}

func configsEqual(a, b *Config) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if len(a.Domains) != len(b.Domains) {
		return false
	}
	for i, domain := range a.Domains {
		if domain != b.Domains[i] {
			return false
		}
	}

	if len(a.Servers) != len(b.Servers) {
		return false
	}
	for i, server := range a.Servers {
		if server != b.Servers[i] {
			return false
		}
	}

	return a.RecordType == b.RecordType &&
		a.Interval == b.Interval &&
		a.AllServers == b.AllServers &&
		a.UntilChange == b.UntilChange &&
		a.OutputFile == b.OutputFile &&
		a.NoColor == b.NoColor &&
		a.ShowHelp == b.ShowHelp &&
		a.ShowVersion == b.ShowVersion
}