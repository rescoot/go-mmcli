# go-mmcli

A Go library for parsing and working with ModemManager CLI (mmcli) JSON output. This library provides a simple interface to work with modem information, connection status, signal strength, location data, network time, and other modem-related functionality.

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
- Get and manage location information
- Establish and manage network connections
- Retrieve network time information

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

### Core Modem Functions
- `Parse(data []byte) (*ModemManager, error)` - Parse mmcli JSON output
- `ListModems() ([]string, error)` - Get full DBus paths of all modems
- `GetModemIDs() ([]string, error)` - Get IDs of all modems
- `GetFirstModemID() (string, error)` - Get ID of first available modem
- `GetModemDetails(id string) (*ModemManager, error)` - Get details for specific modem
- `ResetModem(modemID string) (bool, error)` - Reset a modem
- `GetSIMInfo(simID string) (*SIMInfo, error)` - Get SIM card information

### Modem Information Methods
- `IsConnected() bool` - Check if modem is connected
- `SignalStrength() (int, error)` - Get signal strength percentage
- `GetPortByType(portType string) string` - Get port name by type
- `IsSimLocked() bool` - Check if SIM card is locked
- `GetCurrentAccessTechnology() string` - Get current network technology
- `GetOperatorInfo() (name string, code string)` - Get operator information
- `RemainingUnlockRetries(lockType string) int` - Get remaining unlock attempts
- `IsIPv6Supported() bool` - Check IPv6 support
- `GetAllPorts() map[string]string` - Get all available ports

### Location Functions
- `GetLocationStatus(modemID string) (*LocationStatus, error)` - Get location gathering status
- `GetLocation(modemID string) (*LocationInfo, error)` - Get current location information
- `EnableLocationGathering(modemID string, method string) error` - Enable a location gathering method
- `DisableLocationGathering(modemID string, method string) error` - Disable a location gathering method
- `SetSuplServer(modemID string, address string) error` - Set SUPL server address for A-GPS
- `SetGpsRefreshRate(modemID string, rate int) error` - Set GPS refresh rate
- `EnableLocationSignals(modemID string) error` - Enable location update signaling
- `DisableLocationSignals(modemID string) error` - Disable location update signaling

### Connection Functions
- `Connect(modemID string, settings ConnectSettings) error` - Connect with specific settings
- `Disconnect(modemID string) error` - Disconnect all bearers
- `ConnectWithAPN(modemID string, apn string) error` - Connect with just an APN
- `ConnectWithAuth(modemID string, apn, user, password string) error` - Connect with APN, username and password

### Time Functions
- `GetNetworkTime(modemID string) (*TimeInfo, error)` - Get current network time information
- `GetNetworkTimeAsTime(modemID string) (time.Time, error)` - Get network time as a time.Time object

## Examples

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

### Location Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // Get first available modem ID
    id, err := mmcli.GetFirstModemID()
    if err != nil {
        log.Fatal("No modem found:", err)
    }
    
    // Get location status
    status, err := mmcli.GetLocationStatus(id)
    if err != nil {
        log.Fatal("Failed to get location status:", err)
    }
    
    fmt.Printf("Location capabilities: %v\n", status.Capabilities)
    fmt.Printf("Enabled methods: %v\n", status.Enabled)
    
    // Enable 3GPP location gathering if not already enabled
    for _, method := range status.Capabilities {
        if method == "3gpp-lac-ci" {
            found := false
            for _, enabledMethod := range status.Enabled {
                if enabledMethod == "3gpp-lac-ci" {
                    found = true
                    break
                }
            }
            if !found {
                if err := mmcli.EnableLocationGathering(id, "3gpp"); err != nil {
                    log.Printf("Failed to enable 3GPP location gathering: %v", err)
                }
            }
        }
    }
    
    // Get location information
    location, err := mmcli.GetLocation(id)
    if err != nil {
        log.Fatal("Failed to get location:", err)
    }
    
    fmt.Printf("3GPP Location: MCC=%s, MNC=%s, LAC=%s, CID=%s\n",
        location.ThreeGPP.MCC, location.ThreeGPP.MNC,
        location.ThreeGPP.LAC, location.ThreeGPP.CID)
}
```

### Connection Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // Get first available modem ID
    id, err := mmcli.GetFirstModemID()
    if err != nil {
        log.Fatal("No modem found:", err)
    }
    
    // Connect with APN
    fmt.Println("Connecting to network...")
    err = mmcli.ConnectWithAPN(id, "internet")
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    
    fmt.Println("Connected successfully!")
    
    // Get modem details to check connection
    modem, err := mmcli.GetModemDetails(id)
    if err != nil {
        log.Fatal("Failed to get modem details:", err)
    }
    
    if modem.IsConnected() {
        fmt.Println("Modem is connected")
        
        // Disconnect
        fmt.Println("Disconnecting...")
        err = mmcli.Disconnect(id)
        if err != nil {
            log.Fatal("Failed to disconnect:", err)
        }
        fmt.Println("Disconnected successfully!")
    }
}
```

## Development

### Running Tests

```bash
go test -v
```

### Adding New Features

The library is designed to be easily extended with new ModemManager functionality. To add support for a new mmcli feature:

1. Define appropriate Go structs to represent the data
2. Implement functions that call mmcli with the right arguments
3. Parse the JSON output into the defined structs
4. Add helper methods for common operations
5. Add tests for the new functionality

## License

This project is licensed under [AGPL 3.0](LICENSE).
