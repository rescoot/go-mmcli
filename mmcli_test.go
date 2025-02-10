package mmcli

import (
	"testing"
)

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
