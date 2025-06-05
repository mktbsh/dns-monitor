# DNS Monitor Tool

## Overview
A command-line tool for monitoring DNS record changes and detecting modifications. Efficiently monitors DNS record switching when changing destinations in services like Route 53.

## Functional Requirements

### Core Features
- Monitor specified domain DNS records at regular intervals
- Detect record changes and output logs/notifications
- Simultaneous monitoring of multiple domains
- Simultaneous queries to multiple DNS servers

### Monitoring Operation Modes
1. **Continuous Monitoring Mode** (Default)
   - Persistent checking at specified intervals
   - Continues operation until stopped with Ctrl+C
2. **Monitor Until Change Mode**
   - Terminates process after detecting and notifying changes

### Output & Notifications
- Standard output: Color-coded comparison display of before/after values
- Log file output: Timestamped logs
- macOS notification bar alerts (macOS environment only)

## Technical Specifications

### Supported Environments
- Linux, macOS

### Performance Requirements
- Low memory usage and high efficiency implementation
- Optimized resource usage for multiple domain monitoring

### Implementation Language
- Go (cross-platform support, high efficiency, suitable for CLI tools)

### Implementation Constraints
- **No external libraries allowed**
- Implementation using Go standard library only

## Command Line Specification

### Basic Syntax
```bash
dns-monitor [OPTIONS] DOMAIN [DOMAIN...]
```

### Options
```bash
-t, --type TYPE          DNS record type (A, AAAA, CNAME, MX, TXT, etc.) [default: A]
-i, --interval DURATION  Check interval [default: 5s]
-s, --server SERVER      Specify DNS server (multiple allowed)
--all-servers           Query all major DNS servers (8.8.8.8, 1.1.1.1, 1.0.0.1)
--until-change          Monitor until change mode
-o, --output FILE       Log file output destination
--no-color              Disable colored output
-h, --help              Display help
-v, --version           Display version
```

### Usage Examples
```bash
# Basic monitoring
dns-monitor example.com

# Monitor multiple domains at 30-second intervals
dns-monitor -i 30s example.com api.example.com

# Monitor CNAME record, exit on change detection
dns-monitor -t CNAME --until-change www.example.com

# Simultaneous queries to multiple DNS servers
dns-monitor -s 8.8.8.8 -s 1.1.1.1 example.com

# With log file output
dns-monitor -o /var/log/dns-monitor.log example.com
```

## Output Format

### Standard Output
```
[2025-06-05 15:30:45] example.com (A) - No change: [203.0.113.1]
[2025-06-05 15:30:50] example.com (A) - CHANGE DETECTED:
  Before: [203.0.113.1, 203.0.113.2, 203.0.113.3]
  After:  [203.0.113.1, 203.0.113.4, 203.0.113.5]
[2025-06-05 15:30:50] api.example.com (A) - No change: [198.51.100.1]
```

### Parallel Display for Multiple Domain Monitoring
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

### Color-coded Display
- No change: Green
- Change detected: Red (before) → Blue (after)
- Error: Yellow

## Implementation Details

### Major Components
1. **DNS Client**
   - Parallel queries to each DNS server
   - Multiple IP address sorting processing
   - Timeout handling
2. **Monitoring Engine**
   - Periodic execution scheduler
   - Change detection logic (IP address additions/deletions only)
3. **Output Handler**
   - Color-coded console output
   - Log file writing
   - macOS notifications

### Multiple IP Address Processing Specification
- **Sorting**: Sort retrieved IP addresses in ascending order
- **Change Detection**: Only IP address additions/deletions are treated as changes
- **Order Changes**: Order changes within the same IP address group are not treated as changes
- **Display Format**: Array display in `[IP1, IP2, IP3]` format

### Error Handling
- DNS server connection errors
- Invalid domain names
- Permission errors (log file writing)

### Configuration
- No configuration file support required
- All controlled via command-line options

## Distribution Method
- Single binary
- Distribution via GitHub Releases
- Homebrew installation support (future consideration)

## Development Methodology

### Development Process
Follow this iterative development process for each feature:

1. **Code Generation**
   - Implement the feature according to specifications
   - Write clean, well-documented Go code
   - Follow Go conventions and best practices

2. **Test Code Generation**
   - Create comprehensive unit tests for the implemented feature
   - Include edge cases and error scenarios
   - Use Go's built-in testing framework (`testing` package)
   - Aim for high test coverage

3. **Test Execution**
   - Run all tests using `go test`
   - Ensure all tests pass before proceeding
   - Fix any failing tests immediately

4. **Bug Fixes and Refinement**
   - Address any issues found during testing
   - Refactor code if necessary for better maintainability
   - Re-run tests after modifications

5. **Git Commit**
   - Once tests pass successfully, create a Git commit
   - Use appropriate commit granularity (one feature per commit)
   - Write clear, descriptive commit messages
   - Follow conventional commit format when possible

### Testing Strategy
- **Unit Tests**: Test individual functions and components
- **Integration Tests**: Test DNS query functionality with real servers
- **Error Handling Tests**: Verify proper error handling for network failures
- **Mock Testing**: Use mock DNS responses for consistent testing

### Commit Guidelines
- Commit after each successfully tested feature
- Use Conventional Commit format for commit message titles
- Keep commits atomic and focused on single features
- Include relevant tests in the same commit as the feature
- Use clear, descriptive commit messages without co-author information

#### Commit Message Format
```
<type>(<scope>): <description>

<optional body with detailed explanation>
```

#### Examples
```
feat(dns): add basic DNS monitoring functionality

- Add DNS client for single domain queries
- Implement continuous monitoring loop
- Add basic console output formatting
```

```
test(dns): add unit tests for DNS monitoring

- Add tests for DNS query parsing
- Add tests for change detection logic
- Add error handling test cases
```

```
fix(output): correct IP address sorting logic

- Sort IP addresses before comparison
- Update test cases for sorted output
```

## Development Priority
1. Basic single domain monitoring functionality
2. Multiple domain monitoring
3. Multiple DNS server support
4. macOS notification functionality
5. Advanced output options