# DNS Monitor

[![CI](https://github.com/mktbsh/dns-monitor/workflows/CI/badge.svg)](https://github.com/mktbsh/dns-monitor/actions)
[![Release](https://github.com/mktbsh/dns-monitor/workflows/Release/badge.svg)](https://github.com/mktbsh/dns-monitor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mktbsh/dns-monitor)](https://goreportcard.com/report/github.com/mktbsh/dns-monitor)

A command-line tool for monitoring DNS record changes and detecting modifications. Efficiently monitors DNS record switching when changing destinations in services like Route 53.

## Features

- 🔍 **Real-time DNS Monitoring** - Monitor specified domain DNS records at regular intervals
- 🚨 **Change Detection** - Detect record changes and output logs/notifications
- 🌐 **Multiple Domains** - Simultaneous monitoring of multiple domains
- 🖥️ **Multiple DNS Servers** - Simultaneous queries to multiple DNS servers
- 📋 **Record Type Support** - Support for A, AAAA, CNAME, MX, TXT record types
- 🎨 **Color-coded Output** - Easy-to-read comparison display of before/after values
- 📝 **Logging** - Timestamped log file output
- ⚡ **High Performance** - Low memory usage and efficient implementation
- 🖱️ **Cross-platform** - Linux, macOS, and Windows support

## Installation

### Quick Install

```bash
# Download and install automatically
curl -sSL https://raw.githubusercontent.com/mktbsh/dns-monitor/main/install.sh | bash
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/mktbsh/dns-monitor/releases)
2. Extract the archive for your platform
3. Move the binary to your PATH (e.g., `/usr/local/bin`)
4. Make it executable: `chmod +x dns-monitor`

### Build from Source

```bash
git clone https://github.com/mktbsh/dns-monitor.git
cd dns-monitor
go build -o dns-monitor .
```

## Usage

### Basic Usage

```bash
# Monitor a single domain
dns-monitor example.com

# Monitor multiple domains
dns-monitor example.com api.example.com www.example.com
```

### Advanced Options

```bash
# Monitor with custom interval
dns-monitor -i 30s example.com

# Monitor specific record type
dns-monitor -t CNAME www.example.com

# Monitor until change detected (exit after first change)
dns-monitor --until-change example.com

# Use specific DNS servers
dns-monitor -s 8.8.8.8 -s 1.1.1.1 example.com

# Use all major DNS servers
dns-monitor --all-servers example.com

# Save logs to file
dns-monitor -o /var/log/dns-monitor.log example.com

# Disable colored output
dns-monitor --no-color example.com
```

### Command Line Options

```
OPTIONS:
    -t, --type TYPE          DNS record type (A, AAAA, CNAME, MX, TXT) [default: A]
    -i, --interval DURATION  Check interval [default: 5s]
    -s, --server SERVER      Specify DNS server (multiple allowed)
    --all-servers           Query all major DNS servers (8.8.8.8, 1.1.1.1, 1.0.0.1)
    --until-change          Monitor until change mode
    -o, --output FILE       Log file output destination
    --no-color              Disable colored output
    -h, --help              Display help
    -v, --version           Display version
```

## Output Examples

### Single Domain Monitoring

```
DNS Monitor Tool v1.0.0
Monitoring 1 domain(s) every 5s
Record type: A
Press Ctrl+C to stop

[2025-06-05 15:30:45] example.com (A) - Initial: [203.0.113.1]
[2025-06-05 15:30:50] example.com (A) - No change: [203.0.113.1]
[2025-06-05 15:30:55] example.com (A) - CHANGE DETECTED:
  Before: [203.0.113.1]
  After:  [203.0.113.1, 203.0.113.2]
```

### Multiple Domain Monitoring

```
[2025-06-05 15:30:45]
├─ example.com (A):     [203.0.113.1] (no change)
├─ api.example.com (A): [198.51.100.1, 198.51.100.2] (no change)
└─ www.example.com (A): [203.0.113.1] (no change)

[2025-06-05 15:30:50]
├─ example.com (A):     [203.0.113.1] → [203.0.113.1, 203.0.113.2] (CHANGED)
├─ api.example.com (A): [198.51.100.1, 198.51.100.2] (no change)
└─ www.example.com (A): [203.0.113.1] (no change)
```

### Color Coding

- **Green**: No changes detected or initial records
- **Red**: Changes detected (before values)
- **Blue**: Changes detected (after values)
- **Yellow**: Errors or warnings

## Use Cases

- **DNS Migration Monitoring** - Monitor DNS changes during server migrations
- **Load Balancer Verification** - Verify load balancer DNS updates
- **CDN Monitoring** - Track CDN endpoint changes
- **DNS Propagation Checking** - Monitor DNS propagation across different servers
- **Infrastructure Monitoring** - Alert on unexpected DNS changes

## Technical Details

### Supported Record Types

- **A** - IPv4 addresses
- **AAAA** - IPv6 addresses
- **CNAME** - Canonical name records
- **MX** - Mail exchange records
- **TXT** - Text records

### Change Detection

- Only IP address additions/deletions are treated as changes
- Order changes within the same IP address group are ignored
- All IP addresses are sorted before comparison

### Performance

- Efficient implementation using Go standard library only
- Low memory footprint
- Optimized for monitoring multiple domains simultaneously
- No external dependencies

## Development

### Prerequisites

- Go 1.24 or later

### Building

```bash
# Build for current platform
go build -o dns-monitor .

# Cross-platform build
./build.sh
```

### Testing

```bash
# Run tests
go test -v ./...

# Run tests with race detection
go test -race -v ./...
```

### Local Installation

```bash
# Install from source to custom directory
./install.sh --local --dir ~/.local/bin
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/mktbsh/dns-monitor)
- 🐛 [Issue Tracker](https://github.com/mktbsh/dns-monitor/issues)
- 📦 [Releases](https://github.com/mktbsh/dns-monitor/releases)

---

Made with ❤️ for DNS monitoring needs