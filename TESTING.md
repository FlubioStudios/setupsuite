# Testing Guide for SetupSuite

This document describes the testing strategy and available tests for the SetupSuite CLI tool.

## Test Structure

The project includes comprehensive test coverage at multiple levels:

### Unit Tests

#### Configuration Parser Tests (`suite/config/parser_test.go`)
- **TestParseConfig**: Tests parsing of complete server configuration files
  - Basic web server configuration with SSH, firewall, and tools
  - Database server configuration
  - Minimal configuration with defaults
- **TestParseFirewall**: Tests firewall configuration parsing
  - Multiple ports
  - Single port
- **TestParseInstallTools**: Tests tools installation configuration
  - Multiple tools
  - Single tool

#### Package Manager Tests (`suite/package_manager_test.go`)
- **TestDetectPackageManager**: Tests OS detection and package manager selection
  - Ubuntu (apt)
  - CentOS (yum)
  - Fedora (dnf)
  - Arch (pacman)
  - Alpine (apk)
- **TestInstallPackages**: Tests package installation with different managers
- **TestUpdateSystem**: Tests system update commands

#### Service Manager Tests (`suite/service_manager_test.go`)
- **TestDetectServiceManager**: Tests service manager detection
  - Systemd
  - SysVinit
  - OpenRC
- **TestEnableService**: Tests service enabling
- **TestStartService**: Tests service starting
- **TestServiceStatus**: Tests service status checking

### Integration Tests

#### Full Setup Integration Test (`suite/integration_test.go`)
- **TestFullServerSetup**: Tests complete server setup workflow
  - Config file creation
  - Configuration parsing
  - Package manager detection
  - Service manager detection
  - Setup execution

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Specific Package Tests
```bash
# Config parser tests
go test ./suite/config/

# Package manager tests
go test ./suite/

# Integration tests
go test ./suite/ -run Integration
```

### Run Specific Test
```bash
go test ./suite/config/ -run TestParseConfig
go test ./suite/ -run TestDetectPackageManager
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Data

Test configuration files are located in `testdata/configs/`:
- `test_web.sscfg`: Web server configuration example
- `test_database.sscfg`: Database server configuration example

These files are used by integration tests to validate end-to-end functionality.

## Writing New Tests

### Unit Tests
When adding new functionality, ensure you:
1. Create table-driven tests for multiple scenarios
2. Test both success and error cases
3. Use descriptive test names
4. Add test helpers for common setup

Example:
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NewFeature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
Integration tests should:
1. Test realistic workflows
2. Use actual configuration files from testdata/
3. Verify interactions between components
4. Clean up any created resources

## Continuous Integration

The test suite is designed to run in CI/CD pipelines. Ensure:
- All tests pass before merging
- No tests require root privileges (mock system calls)
- Tests run quickly (< 10 seconds total)
- No external dependencies required

## Test Coverage Goals

- **Unit Tests**: > 80% coverage for all packages
- **Integration Tests**: Cover all main CLI workflows
- **Edge Cases**: Test error handling and boundary conditions

## Mocking

The test suite uses mocking for:
- System command execution
- File system operations
- OS detection
- Package manager/service manager interactions

This allows tests to run without requiring actual system changes or root privileges.
