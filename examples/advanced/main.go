// examples/advanced/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rescoot/go-mmcli"
)

// writeModemDetails writes detailed modem information using a tabwriter for formatting
func writeModemDetails(w *tabwriter.Writer, mm *mmcli.ModemManager) {
	fmt.Fprintln(w, "=== Modem Details ===\t\t")
	fmt.Fprintf(w, "Manufacturer:\t%s\t\n", mm.Modem.Generic.Manufacturer)
	fmt.Fprintf(w, "Model:\t%s\t\n", mm.Modem.Generic.Model)
	fmt.Fprintf(w, "Revision:\t%s\t\n", mm.Modem.Generic.Revision)
	fmt.Fprintf(w, "IMEI:\t%s\t\n", mm.Modem.ThreeGPP.IMEI)
	fmt.Fprintf(w, "Equipment ID:\t%s\t\n", mm.Modem.Generic.EquipmentIdentifier)
	fmt.Fprintln(w, "\t\t")
}

// writeNetworkDetails writes network-related information
func writeNetworkDetails(w *tabwriter.Writer, mm *mmcli.ModemManager) {
	fmt.Fprintln(w, "=== Network Status ===\t\t")

	if mm.IsConnected() {
		strength, _ := mm.SignalStrength()
		fmt.Fprintf(w, "Connection:\tConnected\t\n")
		fmt.Fprintf(w, "Signal:\t%d%%\t\n", strength)

		name, code := mm.GetOperatorInfo()
		fmt.Fprintf(w, "Operator:\t%s (%s)\t\n", name, code)

		tech := mm.GetCurrentAccessTechnology()
		fmt.Fprintf(w, "Technology:\t%s\t\n", tech)

		// Display supported IP families
		ipFamilies := mm.Modem.Generic.SupportedIPFamilies
		fmt.Fprintf(w, "IP Support:\t%s\t\n", strings.Join(ipFamilies, ", "))
	} else {
		fmt.Fprintf(w, "Connection:\tDisconnected\t\n")
	}
	fmt.Fprintln(w, "\t\t")
}

// writePortDetails writes detailed port information
func writePortDetails(w *tabwriter.Writer, mm *mmcli.ModemManager) {
	fmt.Fprintln(w, "=== Port Details ===\t\t")
	ports := mm.GetAllPorts()
	for portType, portName := range ports {
		fmt.Fprintf(w, "%s Port:\t%s\t\n", strings.ToUpper(portType), portName)
	}
	fmt.Fprintln(w, "\t\t")
}

// writeBandInformation writes current and supported band information
func writeBandInformation(w *tabwriter.Writer, mm *mmcli.ModemManager) {
	fmt.Fprintln(w, "=== Band Information ===\t\t")
	fmt.Fprintf(w, "Current Bands:\t%s\t\n", strings.Join(mm.Modem.Generic.CurrentBands, ", "))
	fmt.Fprintf(w, "Supported Bands:\t%s\t\n", strings.Join(mm.Modem.Generic.SupportedBands, ", "))
	fmt.Fprintln(w, "\t\t")
}

// writeModemCapabilities writes detailed capability information
func writeModemCapabilities(w *tabwriter.Writer, mm *mmcli.ModemManager) {
	fmt.Fprintln(w, "=== Modem Capabilities ===\t\t")
	fmt.Fprintf(w, "Current Modes:\t%s\t\n", mm.Modem.Generic.CurrentModes)
	fmt.Fprintf(w, "Current Capabilities:\t%s\t\n", strings.Join(mm.Modem.Generic.CurrentCapabilities, ", "))
	fmt.Fprintf(w, "Supported Capabilities:\t%s\t\n", strings.Join(mm.Modem.Generic.SupportedCapabilities, ", "))
	fmt.Fprintln(w, "\t\t")
}

func main() {
	// Get available modems
	ids, err := mmcli.GetModemIDs()
	if err != nil {
		log.Fatal("Failed to get modem IDs:", err)
	}
	if len(ids) == 0 {
		log.Fatal("No modems found")
	}

	fmt.Printf("Found %d modem(s). Monitoring modem %s...\n", len(ids), ids[0])

	// Create a new tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for {
		// Get modem details
		mm, err := mmcli.GetModemDetails(ids[0])
		if err != nil {
			log.Printf("Failed to get modem details: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Write all sections
		writeModemDetails(w, mm)
		writeNetworkDetails(w, mm)
		writePortDetails(w, mm)
		writeBandInformation(w, mm)
		writeModemCapabilities(w, mm)

		// Flush the tabwriter
		w.Flush()

		// Wait before next update
		time.Sleep(5 * time.Second)
	}
}
