package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Monitor struct {
	config     *Config
	dnsClient  *DNSClient
	lastRecords map[string]*DNSRecord
	logger     *log.Logger
}

func NewMonitor(config *Config) *Monitor {
	dnsClient := NewDNSClient(config.Servers, 5*time.Second)
	
	logger := log.New(os.Stdout, "", 0)
	if config.OutputFile != "" {
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Warning: Failed to open output file %s: %v", config.OutputFile, err)
		} else {
			logger = log.New(file, "", log.LstdFlags)
		}
	}

	return &Monitor{
		config:      config,
		dnsClient:   dnsClient,
		lastRecords: make(map[string]*DNSRecord),
		logger:      logger,
	}
}

func (m *Monitor) Start() error {
	fmt.Printf("DNS Monitor Tool v%s\n", Version)
	fmt.Printf("Monitoring %d domain(s) every %s\n", len(m.config.Domains), m.config.Interval)
	fmt.Printf("Record type: %s\n", m.config.RecordType)
	if len(m.config.Servers) > 0 {
		fmt.Printf("DNS servers: %v\n", m.config.Servers)
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			changed := m.checkDomains()
			if changed && m.config.UntilChange {
				fmt.Println("Change detected. Exiting due to --until-change mode.")
				return nil
			}
		case <-interrupt:
			fmt.Println("\nReceived interrupt signal. Stopping monitor...")
			return nil
		}
	}
}

func (m *Monitor) checkDomains() bool {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	hasChanges := false

	if len(m.config.Domains) == 1 {
		changed := m.checkSingleDomain(m.config.Domains[0], timestamp)
		if changed {
			hasChanges = true
		}
	} else {
		fmt.Printf("[%s]\n", timestamp)
		for i, domain := range m.config.Domains {
			changed := m.checkDomainInGroup(domain, i == len(m.config.Domains)-1)
			if changed {
				hasChanges = true
			}
		}
		fmt.Println()
	}

	return hasChanges
}

func (m *Monitor) checkSingleDomain(domain, timestamp string) bool {
	record, err := m.dnsClient.Query(domain, m.config.RecordType)
	if err != nil {
		message := fmt.Sprintf("[%s] %s (%s) - ERROR: %v", timestamp, domain, m.config.RecordType, err)
		m.printColored(message, ColorYellow)
		m.logger.Println(message)
		return false
	}

	key := fmt.Sprintf("%s:%s", domain, m.config.RecordType)
	lastRecord, exists := m.lastRecords[key]

	if !exists {
		message := fmt.Sprintf("[%s] %s (%s) - Initial: %s", timestamp, domain, m.config.RecordType, record.String())
		m.printColored(message, ColorGreen)
		m.logger.Println(message)
		m.lastRecords[key] = record
		return false
	}

	if !record.Equals(lastRecord) {
		message := fmt.Sprintf("[%s] %s (%s) - CHANGE DETECTED:", timestamp, domain, m.config.RecordType)
		m.printColored(message, ColorRed)
		m.logger.Println(message)

		beforeMsg := fmt.Sprintf("  Before: %s", lastRecord.String())
		afterMsg := fmt.Sprintf("  After:  %s", record.String())
		
		m.printColored(beforeMsg, ColorRed)
		m.printColored(afterMsg, ColorBlue)
		m.logger.Println(beforeMsg)
		m.logger.Println(afterMsg)

		m.lastRecords[key] = record
		return true
	}

	message := fmt.Sprintf("[%s] %s (%s) - No change: %s", timestamp, domain, m.config.RecordType, record.String())
	m.printColored(message, ColorGreen)
	m.logger.Println(message)
	return false
}

func (m *Monitor) checkDomainInGroup(domain string, isLast bool) bool {
	record, err := m.dnsClient.Query(domain, m.config.RecordType)
	if err != nil {
		prefix := "├─"
		if isLast {
			prefix = "└─"
		}
		message := fmt.Sprintf("%s %s (%s): ERROR - %v", prefix, domain, m.config.RecordType, err)
		m.printColored(message, ColorYellow)
		m.logger.Printf("ERROR: %s (%s) - %v", domain, m.config.RecordType, err)
		return false
	}

	key := fmt.Sprintf("%s:%s", domain, m.config.RecordType)
	lastRecord, exists := m.lastRecords[key]

	prefix := "├─"
	if isLast {
		prefix = "└─"
	}

	if !exists {
		message := fmt.Sprintf("%s %s (%s): %s (initial)", prefix, domain, m.config.RecordType, record.String())
		m.printColored(message, ColorGreen)
		m.logger.Printf("INITIAL: %s (%s) - %s", domain, m.config.RecordType, record.String())
		m.lastRecords[key] = record
		return false
	}

	if !record.Equals(lastRecord) {
		message := fmt.Sprintf("%s %s (%s): %s → %s (CHANGED)", prefix, domain, m.config.RecordType, lastRecord.String(), record.String())
		m.printColored(message, ColorRed)
		m.logger.Printf("CHANGE: %s (%s) - %s → %s", domain, m.config.RecordType, lastRecord.String(), record.String())
		m.lastRecords[key] = record
		return true
	}

	message := fmt.Sprintf("%s %s (%s): %s (no change)", prefix, domain, m.config.RecordType, record.String())
	m.printColored(message, ColorGreen)
	return false
}

func (m *Monitor) printColored(message string, color string) {
	if m.config.NoColor {
		fmt.Println(message)
	} else {
		fmt.Printf("%s%s%s\n", color, message, ColorReset)
	}
}