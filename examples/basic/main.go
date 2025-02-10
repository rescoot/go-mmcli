package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/rescoot/go-mmcli"
)

func main() {
	// Simple monitoring loop to demonstrate basic modem information
	for {
		// Run mmcli command to get modem information
		cmd := exec.Command("mmcli", "-m", "0", "-J")
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Failed to run mmcli: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Parse the output
		mm, err := mmcli.Parse(output)
		if err != nil {
			log.Printf("Failed to parse mmcli output: %v", err)
			// Pretty print the JSON to help debug parsing issues
			var prettyJSON map[string]interface{}
			if err := json.Unmarshal(output, &prettyJSON); err == nil {
				prettyOutput, _ := json.MarshalIndent(prettyJSON, "", "  ")
				log.Printf("Raw JSON: %s", string(prettyOutput))
			}
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
