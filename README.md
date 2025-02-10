# go-mmcli

A Go library for parsing and working with ModemManager CLI (mmcli) JSON output. This library provides a simple interface to work with modem information, connection status, signal strength, and other modem-related functionality.

## Installation

```bash
go get github.com/rescoot/go-mmcli
```

## Features

- Parse mmcli JSON output into Go structs
- Get modem connection status
- Monitor signal strength
- Access port information (QMI, AT, NET, etc.)
- Check SIM lock status
- Get current network technology (2G/3G/4G)
- Get operator information
- Monitor unlock attempts
- Check IPv6 support
- Access complete port mapping

## Finding Modems

The library provides several ways to find modems:

```go
// Get all modem IDs
ids, err := mmcli.GetModemIDs()

// Get full DBus paths
paths, err := mmcli.ListModems()

// Get just the first modem ID (common case)
id, err := mmcli.GetFirstModemID()

// Get details for a specific modem
modem, err := mmcli.GetModemDetails(id)
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // List available modems
    ids, err := mmcli.GetModemIDs()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d modems: %v\n", len(ids), ids)
    
    // Get details for first modem
    modem, err := mmcli.GetModemDetails(ids[0])
    if err != nil {
        log.Fatal(err)
    }
    
    if modem.IsConnected() {
        strength, _ := modem.SignalStrength()
        name, code := modem.GetOperatorInfo()
        tech := modem.GetCurrentAccessTechnology()
        
        fmt.Printf("Connected to %s (%s)\n", name, code)
        fmt.Printf("Technology: %s\n", tech)
        fmt.Printf("Signal: %d%%\n", strength)
    }
}
```

## Available Methods

- `Parse(data []byte) (*ModemManager, error)` - Parse mmcli JSON output
- `ListModems() ([]string, error)` - Get full DBus paths of all modems
- `GetModemIDs() ([]string, error)` - Get IDs of all modems
- `GetFirstModemID() (string, error)` - Get ID of first available modem
- `GetModemDetails(id string) (*ModemManager, error)` - Get details for specific modem
- `IsConnected() bool` - Check if modem is connected
- `SignalStrength() (int, error)` - Get signal strength percentage
- `GetPortByType(portType string) string` - Get port name by type
- `IsSimLocked() bool` - Check if SIM card is locked
- `GetCurrentAccessTechnology() string` - Get current network technology
- `GetOperatorInfo() (name string, code string)` - Get operator information
- `RemainingUnlockRetries(lockType string) int` - Get remaining unlock attempts
- `IsIPv6Supported() bool` - Check IPv6 support
- `GetAllPorts() map[string]string` - Get all available ports

## Development

### Running Tests

```bash
go test -v
```

## License

This project is licensed under [AGPL 3.0](LICENSE).
