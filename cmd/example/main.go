// cmd/example/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rescoot/go-mmcli"
)

var (
	interval = flag.Duration("interval", 5*time.Second, "Update interval")
	mode     = flag.String("mode", "basic", "Display mode (basic, network, ports, all)")
	listOnly = flag.Bool("list", false, "List available modems and exit")
	modemID  = flag.String("modem", "", "Modem ID to monitor (default: first available modem)")
)

func displayBasicInfo(mm *mmcli.ModemManager) {
	if mm.IsConnected() {
		fmt.Println("Status: Connected")

		strength, err := mm.SignalStrength()
		if err == nil {
			fmt.Printf("Signal Strength: %d%%\n", strength)
		}

		name, code := mm.GetOperatorInfo()
		fmt.Printf("Operator: %s (%s)\n", name, code)

		tech := mm.GetCurrentAccessTechnology()
		fmt.Printf("Network Type: %s\n", tech)
	} else {
		fmt.Println("Status: Not Connected")
	}

	if mm.IsSimLocked() {
		fmt.Println("SIM Status: Locked")
		retries := mm.RemainingUnlockRetries("sim-pin")
		fmt.Printf("Remaining PIN attempts: %d\n", retries)
	} else {
		fmt.Println("SIM Status: Unlocked")
	}
}

func displayNetworkInfo(mm *mmcli.ModemManager) {
	fmt.Printf("Registration: %s\n", mm.Modem.ThreeGPP.RegistrationState)
	fmt.Printf("Operator: %s (%s)\n", mm.Modem.ThreeGPP.OperatorName, mm.Modem.ThreeGPP.OperatorCode)
	fmt.Printf("Access Tech: %s\n", mm.GetCurrentAccessTechnology())
	fmt.Printf("Current Bands: %s\n", strings.Join(mm.Modem.Generic.CurrentBands, ", "))
	fmt.Printf("IP Families: %s\n", strings.Join(mm.Modem.Generic.SupportedIPFamilies, ", "))
}

func displayPortInfo(mm *mmcli.ModemManager) {
	fmt.Println("Available Ports:")
	ports := mm.GetAllPorts()
	for portType, portName := range ports {
		fmt.Printf("- %s: %s\n", strings.ToUpper(portType), portName)
	}
}

func main() {
	flag.Parse()

	// List modems if requested
	if *listOnly {
		paths, err := mmcli.ListModems()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Available modems:")
		for _, path := range paths {
			fmt.Printf("  %s\n", path)
		}
		return
	}

	// Get modem ID to monitor
	id := *modemID
	if id == "" {
		var err error
		id, err = mmcli.GetFirstModemID()
		if err != nil {
			log.Fatal("No modem found:", err)
		}
	}

	fmt.Printf("Monitoring modem %s...\n", id)

	for {
		// Get modem details
		mm, err := mmcli.GetModemDetails(id)
		if err != nil {
			log.Printf("Failed to get modem details: %v", err)
			time.Sleep(*interval)
			continue
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Display timestamp
		fmt.Printf("=== Modem %s Status (Updated: %s) ===\n\n",
			id, time.Now().Format("15:04:05"))

		// Display information based on mode
		switch strings.ToLower(*mode) {
		case "basic":
			displayBasicInfo(mm)
		case "network":
			displayNetworkInfo(mm)
		case "ports":
			displayPortInfo(mm)
		case "all":
			fmt.Println("=== Basic Information ===")
			displayBasicInfo(mm)
			fmt.Println("\n=== Network Information ===")
			displayNetworkInfo(mm)
			fmt.Println("\n=== Port Information ===")
			displayPortInfo(mm)
		default:
			fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", *mode)
			os.Exit(1)
		}

		time.Sleep(*interval)
	}
}
