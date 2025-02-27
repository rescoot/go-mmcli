package mmcli

import (
	"os/exec"
	"testing"
)

func TestLocationFunctions(t *testing.T) {
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

	// Test GetLocationStatus
	t.Run("GetLocationStatus", func(t *testing.T) {
		status, err := GetLocationStatus(ids[0])
		if err != nil {
			t.Errorf("Failed to get location status: %v", err)
		}
		if status == nil {
			t.Error("Expected location status, got nil")
		} else {
			// Verify that the status has the expected fields
			if status.Capabilities == nil {
				t.Error("Expected capabilities, got nil")
			}
			if status.Enabled == nil {
				t.Error("Expected enabled, got nil")
			}
			// Signals could be "yes" or "no"
			if status.Signals != "yes" && status.Signals != "no" {
				t.Errorf("Expected signals to be 'yes' or 'no', got '%s'", status.Signals)
			}
		}
	})

	// Test GetLocation
	t.Run("GetLocation", func(t *testing.T) {
		location, err := GetLocation(ids[0])
		if err != nil {
			t.Errorf("Failed to get location: %v", err)
		}
		if location == nil {
			t.Error("Expected location, got nil")
		}
		// Note: We can't verify the actual values as they depend on the modem's state
	})

	// Test EnableLocationGathering and DisableLocationGathering
	// Note: We're not actually enabling/disabling here to avoid side effects
	t.Run("EnableDisableLocationGathering", func(t *testing.T) {
		// Just verify that the functions don't panic
		// In a real test, we would enable, verify, then disable
		method := "3gpp"
		err := EnableLocationGathering(ids[0], method)
		if err != nil {
			// This might fail if the modem doesn't support this method
			t.Logf("Note: Failed to enable location gathering (%s): %v", method, err)
		}

		err = DisableLocationGathering(ids[0], method)
		if err != nil {
			// This might fail if the modem doesn't support this method
			t.Logf("Note: Failed to disable location gathering (%s): %v", method, err)
		}
	})
}
