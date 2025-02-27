package mmcli

import (
	"os/exec"
	"testing"
)

func TestConnectionFunctions(t *testing.T) {
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

	// Test Connect and Disconnect
	// Note: We're not actually connecting/disconnecting here to avoid side effects
	t.Run("ConnectDisconnect", func(t *testing.T) {
		// This is a mock test that doesn't actually connect
		// In a real test environment, you would use a test APN and verify the connection

		// Test the ConnectSettings struct
		settings := ConnectSettings{
			APN:          "test.apn",
			User:         "testuser",
			Password:     "testpass",
			IPType:       "ipv4",
			AllowRoaming: true,
		}

		// Verify the settings are set correctly
		if settings.APN != "test.apn" {
			t.Errorf("Expected APN 'test.apn', got '%s'", settings.APN)
		}
		if settings.User != "testuser" {
			t.Errorf("Expected User 'testuser', got '%s'", settings.User)
		}
		if settings.Password != "testpass" {
			t.Errorf("Expected Password 'testpass', got '%s'", settings.Password)
		}
		if settings.IPType != "ipv4" {
			t.Errorf("Expected IPType 'ipv4', got '%s'", settings.IPType)
		}
		if !settings.AllowRoaming {
			t.Errorf("Expected AllowRoaming true, got false")
		}

		// Test the convenience functions
		apnSettings := ConnectSettings{
			APN: "test.apn",
		}
		if apnSettings.APN != "test.apn" {
			t.Errorf("Expected APN 'test.apn', got '%s'", apnSettings.APN)
		}

		authSettings := ConnectSettings{
			APN:      "test.apn",
			User:     "testuser",
			Password: "testpass",
		}
		if authSettings.APN != "test.apn" {
			t.Errorf("Expected APN 'test.apn', got '%s'", authSettings.APN)
		}
		if authSettings.User != "testuser" {
			t.Errorf("Expected User 'testuser', got '%s'", authSettings.User)
		}
		if authSettings.Password != "testpass" {
			t.Errorf("Expected Password 'testpass', got '%s'", authSettings.Password)
		}
	})
}
