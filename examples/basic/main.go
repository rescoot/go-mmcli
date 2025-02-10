// examples/basic/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rescoot/go-mmcli"
)

func main() {
	// Get first available modem ID
	id, err := mmcli.GetFirstModemID()
	if err != nil {
		log.Fatal("No modem found:", err)
	}

	fmt.Printf("Monitoring modem %s...\n", id)

	// Simple monitoring loop to demonstrate basic modem information
	for {
		// Get modem details
		mm, err := mmcli.GetModemDetails(id)
		if err != nil {
			log.Printf("Failed to get modem details: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Clear screen (simple terminal output refresh)
		fmt.Print("\033[H\033[2J")

		// Display basic modem information
		fmt.Println("=== Basic Modem Status ===")
		if mm.IsConnected() {
			fmt.Println("Status: Connected")

			// Get and display signal strength
			if strength, err := mm.SignalStrength(); err == nil {
				fmt.Printf("Signal Strength: %d%%\n", strength)
			}

			// Get and display operator information
			name, code := mm.GetOperatorInfo()
			fmt.Printf("Operator: %s (%s)\n", name, code)

			// Display current network technology
			tech := mm.GetCurrentAccessTechnology()
			fmt.Printf("Network Type: %s\n", tech)
		} else {
			fmt.Println("Status: Not Connected")
		}

		// Check SIM status
		if mm.IsSimLocked() {
			fmt.Println("\nSIM Status: Locked")
			retries := mm.RemainingUnlockRetries("sim-pin")
			fmt.Printf("Remaining PIN attempts: %d\n", retries)
		} else {
			fmt.Println("\nSIM Status: Unlocked")
		}

		// Sleep for a bit before next update
		time.Sleep(5 * time.Second)
	}
}
