package mmcli

import (
	"encoding/json"
	"os/exec"
	"testing"
)

func TestParseModemList(t *testing.T) {
	jsonData := []byte(`{
		"modem-manager": {
			"version": "1.18.4",
			"modems": [
				"/org/freedesktop/ModemManager1/Modem/0",
				"/org/freedesktop/ModemManager1/Modem/1",
				"/org/freedesktop/ModemManager1/Modem/2"
			]
		}
	}`)

	var list ModemList
	err := json.Unmarshal(jsonData, &list)
	if err != nil {
		t.Fatalf("Failed to parse modem list JSON: %v", err)
	}

	if len(list.ModemManager.Modems) != 3 {
		t.Errorf("Expected 3 modems, got %d", len(list.ModemManager.Modems))
	}

	if list.ModemManager.Version != "1.18.4" {
		t.Errorf("Expected version 1.18.4, got %s", list.ModemManager.Version)
	}

	expectedPaths := []string{
		"/org/freedesktop/ModemManager1/Modem/0",
		"/org/freedesktop/ModemManager1/Modem/1",
		"/org/freedesktop/ModemManager1/Modem/2",
	}

	for i, path := range list.ModemManager.Modems {
		if path != expectedPaths[i] {
			t.Errorf("Expected path %s, got %s", expectedPaths[i], path)
		}
	}
}

func TestGetModemDetails(t *testing.T) {
	// Skip if we're not in an environment with mmcli available
	if _, err := exec.LookPath("mmcli"); err != nil {
		t.Skip("mmcli not available, skipping test")
	}

	// First get available modems
	ids, err := GetModemIDs()
	if err != nil {
		t.Skip("Could not get modem IDs:", err)
	}
	if len(ids) == 0 {
		t.Skip("No modems available for testing")
	}

	// Try to get details for the first modem
	mm, err := GetModemDetails(ids[0])
	if err != nil {
		t.Errorf("Failed to get modem details: %v", err)
	}
	if mm == nil {
		t.Error("Expected modem details, got nil")
	}
}

func TestParse(t *testing.T) {
	jsonData := []byte(`{
		"modem": {
			"3gpp": {
				"imei": "123456789010213",
				"operator-code": "26201",
				"operator-name": "TDG",
				"registration-state": "home"
			},
			"generic": {
				"state": "connected",
				"signal-quality": {
					"recent": "yes",
					"value": "78"
				},
				"ports": [
					"cdc-wdm0 (qmi)",
					"ttyUSB2 (at)",
					"wwan0 (net)"
				],
				"supported-ip-families": [
					"ipv4",
					"ipv6",
					"ipv4v6"
				],
				"access-technologies": [
					"lte"
				],
				"unlock-required": "sim-pin",
				"unlock-retries": [
					"sim-pin (3)",
					"sim-puk (10)"
				]
			}
		}
	}`)

	mm, err := Parse(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Test IsConnected
	if !mm.IsConnected() {
		t.Error("Expected modem to be connected")
	}

	// Test SignalStrength
	strength, err := mm.SignalStrength()
	if err != nil {
		t.Errorf("Failed to get signal strength: %v", err)
	}
	if strength != 78 {
		t.Errorf("Expected signal strength 78, got %d", strength)
	}

	// Test GetPortByType
	if port := mm.GetPortByType("qmi"); port != "cdc-wdm0" {
		t.Errorf("Expected QMI port cdc-wdm0, got %s", port)
	}

	// Test IsSimLocked
	if !mm.IsSimLocked() {
		t.Error("Expected SIM to be locked")
	}

	// Test GetCurrentAccessTechnology
	if tech := mm.GetCurrentAccessTechnology(); tech != "4G" {
		t.Errorf("Expected 4G technology, got %s", tech)
	}

	// Test GetOperatorInfo
	name, code := mm.GetOperatorInfo()
	if name != "TDG" || code != "26201" {
		t.Errorf("Expected operator TDG/26201, got %s/%s", name, code)
	}

	// Test RemainingUnlockRetries
	if retries := mm.RemainingUnlockRetries("sim-pin"); retries != 3 {
		t.Errorf("Expected 3 PIN retries, got %d", retries)
	}

	// Test IsIPv6Supported
	if !mm.IsIPv6Supported() {
		t.Error("Expected IPv6 to be supported")
	}

	// Test GetAllPorts
	ports := mm.GetAllPorts()
	if len(ports) != 3 {
		t.Errorf("Expected 3 ports, got %d", len(ports))
	}
	if ports["qmi"] != "cdc-wdm0" {
		t.Errorf("Expected QMI port cdc-wdm0, got %s", ports["qmi"])
	}
}
