package main

import (
	"testing"
	"time"
)

func TestNewDNSClient(t *testing.T) {
	tests := []struct {
		name     string
		servers  []string
		expected []string
	}{
		{
			name:     "empty servers should use default",
			servers:  []string{},
			expected: []string{"8.8.8.8:53"},
		},
		{
			name:     "custom servers should be preserved",
			servers:  []string{"1.1.1.1:53", "8.8.8.8:53"},
			expected: []string{"1.1.1.1:53", "8.8.8.8:53"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewDNSClient(tt.servers, 5*time.Second)
			if len(client.servers) != len(tt.expected) {
				t.Errorf("expected %d servers, got %d", len(tt.expected), len(client.servers))
			}
			for i, server := range client.servers {
				if server != tt.expected[i] {
					t.Errorf("expected server %s, got %s", tt.expected[i], server)
				}
			}
		})
	}
}

func TestDNSRecord_String(t *testing.T) {
	tests := []struct {
		name     string
		record   *DNSRecord
		expected string
	}{
		{
			name: "single value",
			record: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			expected: "[192.168.1.1]",
		},
		{
			name: "multiple values",
			record: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1", "192.168.1.2"},
			},
			expected: "[192.168.1.1, 192.168.1.2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.record.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDNSRecord_Equals(t *testing.T) {
	tests := []struct {
		name     string
		record1  *DNSRecord
		record2  *DNSRecord
		expected bool
	}{
		{
			name: "identical records",
			record1: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			record2: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			expected: true,
		},
		{
			name: "different domains",
			record1: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			record2: &DNSRecord{
				Domain: "test.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			expected: false,
		},
		{
			name: "different types",
			record1: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			record2: &DNSRecord{
				Domain: "example.com",
				Type:   "AAAA",
				Values: []string{"192.168.1.1"},
			},
			expected: false,
		},
		{
			name: "different values",
			record1: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			record2: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.2"},
			},
			expected: false,
		},
		{
			name: "different number of values",
			record1: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1"},
			},
			record2: &DNSRecord{
				Domain: "example.com",
				Type:   "A",
				Values: []string{"192.168.1.1", "192.168.1.2"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.record1.Equals(tt.record2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDNSClient_Query_UnsupportedType(t *testing.T) {
	client := NewDNSClient([]string{}, 5*time.Second)
	_, err := client.Query("example.com", "UNSUPPORTED")
	if err == nil {
		t.Error("expected error for unsupported record type")
	}
	expectedError := "unsupported record type: UNSUPPORTED"
	if err.Error() != expectedError {
		t.Errorf("expected error %s, got %s", expectedError, err.Error())
	}
}