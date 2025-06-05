package main

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	config := &Config{
		Domains:    []string{"example.com"},
		RecordType: "A",
		Interval:   5 * time.Second,
		Servers:    []string{"8.8.8.8:53"},
	}

	monitor := NewMonitor(config)

	if monitor.config != config {
		t.Error("Monitor config not set correctly")
	}

	if monitor.dnsClient == nil {
		t.Error("DNS client not initialized")
	}

	if monitor.lastRecords == nil {
		t.Error("Last records map not initialized")
	}

	if monitor.logger == nil {
		t.Error("Logger not initialized")
	}
}

func TestMonitor_printColored(t *testing.T) {
	tests := []struct {
		name     string
		noColor  bool
		message  string
		color    string
		expected string
	}{
		{
			name:     "with color",
			noColor:  false,
			message:  "test message",
			color:    ColorRed,
			expected: "\033[31mtest message\033[0m",
		},
		{
			name:     "no color",
			noColor:  true,
			message:  "test message",
			color:    ColorRed,
			expected: "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				NoColor: tt.noColor,
			}

			var buf bytes.Buffer
			monitor := &Monitor{
				config: config,
				logger: log.New(&buf, "", 0),
			}

			var output bytes.Buffer
			oldOut := log.Writer()
			log.SetOutput(&output)
			defer log.SetOutput(oldOut)

			monitor.printColored(tt.message, tt.color)
		})
	}
}

func TestMonitor_checkDomainInGroup(t *testing.T) {
	config := &Config{
		Domains:    []string{"example.com"},
		RecordType: "A",
		Interval:   5 * time.Second,
		NoColor:    true,
	}

	var buf bytes.Buffer
	monitor := &Monitor{
		config:      config,
		dnsClient:   NewDNSClient([]string{}, 5*time.Second),
		lastRecords: make(map[string]*DNSRecord),
		logger:      log.New(&buf, "", 0),
	}

	record1 := &DNSRecord{
		Domain: "example.com",
		Type:   "A",
		Values: []string{"192.168.1.1"},
	}


	key := "example.com:A"
	monitor.lastRecords[key] = record1

	changed := monitor.checkDomainInGroup("example.com", true)

	if !changed {
		t.Error("Expected change detection when record is different")
	}

	newRecord := monitor.lastRecords[key]
	if newRecord.Equals(record1) {
		t.Error("Last record should have been updated")
	}
}

func TestColors(t *testing.T) {
	if ColorReset != "\033[0m" {
		t.Errorf("ColorReset should be \\033[0m, got %s", ColorReset)
	}
	if ColorRed != "\033[31m" {
		t.Errorf("ColorRed should be \\033[31m, got %s", ColorRed)
	}
	if ColorGreen != "\033[32m" {
		t.Errorf("ColorGreen should be \\033[32m, got %s", ColorGreen)
	}
	if ColorYellow != "\033[33m" {
		t.Errorf("ColorYellow should be \\033[33m, got %s", ColorYellow)
	}
	if ColorBlue != "\033[34m" {
		t.Errorf("ColorBlue should be \\033[34m, got %s", ColorBlue)
	}
}