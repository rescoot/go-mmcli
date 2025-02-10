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

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    "os/exec"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // Get JSON output from mmcli
    cmd := exec.Command("mmcli", "-m", "0", "-J")
    output, err := cmd.Output()
    if err != nil {
        log.Fatal("Failed to run mmcli:", err)
    }
    
    // Parse the output
    mm, err := mmcli.Parse(output)
    if err != nil {
        log.Fatal("Failed to parse mmcli output:", err)
    }
    
    // Check connection status
    if mm.IsConnected() {
        fmt.Println("Modem is connected!")
        
        // Get signal strength
        strength, err := mm.SignalStrength()
        if err == nil {
            fmt.Printf("Signal strength: %d%%\n", strength)
        }
        
        // Get operator info
        name, code := mm.GetOperatorInfo()
        fmt.Printf("Connected to: %s (Code: %s)\n", name, code)
        
        // Get network technology
        tech := mm.GetCurrentAccessTechnology()
        fmt.Printf("Using %s technology\n", tech)
    } else {
        fmt.Println("Modem is not connected")
    }
}
```

### Working with Ports

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // Parse mmcli output...
    mm, err := mmcli.Parse(output)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get all ports
    ports := mm.GetAllPorts()
    fmt.Println("Available ports:")
    for portType, portName := range ports {
        fmt.Printf("- %s: %s\n", portType, portName)
    }
    
    // Get specific port types
    if qmiPort := mm.GetPortByType("qmi"); qmiPort != "" {
        fmt.Printf("QMI port: %s\n", qmiPort)
    }
    
    if atPort := mm.GetPortByType("at"); atPort != "" {
        fmt.Printf("AT command port: %s\n", atPort)
    }
    
    if netPort := mm.GetPortByType("net"); netPort != "" {
        fmt.Printf("Network interface: %s\n", netPort)
    }
}
```

### Checking SIM Status

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rescoot/go-mmcli"
)

func main() {
    // Parse mmcli output...
    mm, err := mmcli.Parse(output)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if SIM is locked
    if mm.IsSimLocked() {
        fmt.Println("SIM card is locked!")
        
        // Check remaining PIN attempts
        pinRetries := mm.RemainingUnlockRetries("sim-pin")
        fmt.Printf("Remaining PIN attempts: %d\n", pinRetries)
        
        pukRetries := mm.RemainingUnlockRetries("sim-puk")
        fmt.Printf("Remaining PUK attempts: %d\n", pukRetries)
    } else {
        fmt.Println("SIM card is unlocked")
    }
}
```

## Available Methods

- `Parse(data []byte) (*ModemManager, error)` - Parse mmcli JSON output
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
