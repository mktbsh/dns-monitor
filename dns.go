package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

type DNSClient struct {
	servers []string
	timeout time.Duration
}

type DNSRecord struct {
	Domain string
	Type   string
	Values []string
}

func NewDNSClient(servers []string, timeout time.Duration) *DNSClient {
	if len(servers) == 0 {
		servers = []string{"8.8.8.8:53"}
	}
	return &DNSClient{
		servers: servers,
		timeout: timeout,
	}
}

func (c *DNSClient) Query(domain, recordType string) (*DNSRecord, error) {
	switch strings.ToUpper(recordType) {
	case "A":
		return c.queryA(domain)
	case "AAAA":
		return c.queryAAAA(domain)
	case "CNAME":
		return c.queryCNAME(domain)
	case "MX":
		return c.queryMX(domain)
	case "TXT":
		return c.queryTXT(domain)
	default:
		return nil, fmt.Errorf("unsupported record type: %s", recordType)
	}
}

func (c *DNSClient) queryA(domain string) (*DNSRecord, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup A record for %s: %v", domain, err)
	}

	var ipv4s []string
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipv4s = append(ipv4s, ip.String())
		}
	}

	if len(ipv4s) == 0 {
		return nil, fmt.Errorf("no A records found for %s", domain)
	}

	sort.Strings(ipv4s)
	return &DNSRecord{
		Domain: domain,
		Type:   "A",
		Values: ipv4s,
	}, nil
}

func (c *DNSClient) queryAAAA(domain string) (*DNSRecord, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup AAAA record for %s: %v", domain, err)
	}

	var ipv6s []string
	for _, ip := range ips {
		if ip.To4() == nil && ip.To16() != nil {
			ipv6s = append(ipv6s, ip.String())
		}
	}

	if len(ipv6s) == 0 {
		return nil, fmt.Errorf("no AAAA records found for %s", domain)
	}

	sort.Strings(ipv6s)
	return &DNSRecord{
		Domain: domain,
		Type:   "AAAA",
		Values: ipv6s,
	}, nil
}

func (c *DNSClient) queryCNAME(domain string) (*DNSRecord, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup CNAME record for %s: %v", domain, err)
	}

	cname = strings.TrimSuffix(cname, ".")
	return &DNSRecord{
		Domain: domain,
		Type:   "CNAME",
		Values: []string{cname},
	}, nil
}

func (c *DNSClient) queryMX(domain string) (*DNSRecord, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup MX record for %s: %v", domain, err)
	}

	var mxValues []string
	for _, mx := range mxRecords {
		mxValues = append(mxValues, fmt.Sprintf("%d %s", mx.Pref, strings.TrimSuffix(mx.Host, ".")))
	}

	sort.Strings(mxValues)
	return &DNSRecord{
		Domain: domain,
		Type:   "MX",
		Values: mxValues,
	}, nil
}

func (c *DNSClient) queryTXT(domain string) (*DNSRecord, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup TXT record for %s: %v", domain, err)
	}

	if len(txtRecords) == 0 {
		return nil, fmt.Errorf("no TXT records found for %s", domain)
	}

	sort.Strings(txtRecords)
	return &DNSRecord{
		Domain: domain,
		Type:   "TXT",
		Values: txtRecords,
	}, nil
}

func (r *DNSRecord) String() string {
	return fmt.Sprintf("[%s]", strings.Join(r.Values, ", "))
}

func (r *DNSRecord) Equals(other *DNSRecord) bool {
	if r.Domain != other.Domain || r.Type != other.Type {
		return false
	}
	if len(r.Values) != len(other.Values) {
		return false
	}
	for i, v := range r.Values {
		if v != other.Values[i] {
			return false
		}
	}
	return true
}