package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rescoot/go-mmcli"
)

var (
	modemIndex = flag.String("modem", "0", "Modem index to monitor")
	interval   = flag.Duration("interval", 5*time.Second, "Update interval")
	mode       = flag.String("mode", "basic", "Display mode (basic, network, ports, all)")
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

	for {
		// Run mmcli command
		cmd := exec.Command("mmcli", "-m", *modemIndex, "-J")
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Failed to run mmcli: %v", err)
			time.Sleep(*interval)
			continue
		}

		// Parse the output
		mm, err := mmcli.Parse(output)
		if err != nil {
			log.Printf("Failed to parse mmcli output: %v", err)
			time.Sleep(*interval)
			continue
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Display timestamp
		fmt.Printf("=== Modem Status (Updated: %s) ===\n\n", time.Now().Format("15:04:05"))

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
