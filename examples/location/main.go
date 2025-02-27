// examples/location/main.go
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

	fmt.Printf("Using modem %s\n\n", id)

	// Get location status
	fmt.Println("=== Location Status ===")
	status, err := mmcli.GetLocationStatus(id)
	if err != nil {
		log.Printf("Failed to get location status: %v", err)
	} else {
		fmt.Printf("Capabilities: %v\n", status.Capabilities)
		fmt.Printf("Enabled: %v\n", status.Enabled)
		fmt.Printf("Signals: %s\n", status.Signals)
		fmt.Printf("GPS Refresh Rate: %s\n", status.GPS.RefreshRate)
		fmt.Printf("SUPL Server: %s\n", status.GPS.SuplServer)
	}

	// Enable location gathering if not already enabled
	fmt.Println("\n=== Enabling Location Gathering ===")
	enabled := false
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
				fmt.Println("Enabling 3GPP location gathering...")
				if err := mmcli.EnableLocationGathering(id, "3gpp"); err != nil {
					log.Printf("Failed to enable 3GPP location gathering: %v", err)
				} else {
					enabled = true
				}
			} else {
				fmt.Println("3GPP location gathering already enabled")
				enabled = true
			}
		}
	}

	if enabled {
		// Get location information
		fmt.Println("\n=== Location Information ===")
		location, err := mmcli.GetLocation(id)
		if err != nil {
			log.Printf("Failed to get location: %v", err)
		} else {
			fmt.Printf("3GPP Location:\n")
			fmt.Printf("  MCC: %s\n", location.ThreeGPP.MCC)
			fmt.Printf("  MNC: %s\n", location.ThreeGPP.MNC)
			fmt.Printf("  LAC: %s\n", location.ThreeGPP.LAC)
			fmt.Printf("  CID: %s\n", location.ThreeGPP.CID)
			fmt.Printf("  TAC: %s\n", location.ThreeGPP.TAC)

			fmt.Printf("\nGPS Location:\n")
			fmt.Printf("  Latitude: %s\n", location.GPS.Latitude)
			fmt.Printf("  Longitude: %s\n", location.GPS.Longitude)
			fmt.Printf("  Altitude: %s\n", location.GPS.Altitude)
			fmt.Printf("  UTC: %s\n", location.GPS.UTC)

			fmt.Printf("\nCDMA Location:\n")
			fmt.Printf("  Latitude: %s\n", location.CDMABS.Latitude)
			fmt.Printf("  Longitude: %s\n", location.CDMABS.Longitude)
		}
	}

	// Try to get network time
	fmt.Println("\n=== Network Time ===")
	timeInfo, err := mmcli.GetNetworkTime(id)
	if err != nil {
		log.Printf("Failed to get network time: %v", err)
	} else {
		fmt.Printf("Network Time: %s\n", timeInfo.NetworkTime)
		fmt.Printf("Local Time: %s\n", timeInfo.Local)
	}

	// Try to connect (this is just an example, replace with your actual APN)
	fmt.Println("\n=== Connection Management ===")
	fmt.Println("Connecting to network...")
	err = mmcli.ConnectWithAPN(id, "internet")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	} else {
		fmt.Println("Connected successfully!")

		// Wait a bit before disconnecting
		fmt.Println("Waiting 5 seconds before disconnecting...")
		time.Sleep(5 * time.Second)

		// Disconnect
		fmt.Println("Disconnecting...")
		err = mmcli.Disconnect(id)
		if err != nil {
			log.Printf("Failed to disconnect: %v", err)
		} else {
			fmt.Println("Disconnected successfully!")
		}
	}
}
