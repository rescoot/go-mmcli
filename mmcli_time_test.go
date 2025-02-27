package mmcli

import (
	"os/exec"
	"testing"
	"time"
)

func TestTimeFunctions(t *testing.T) {
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

	// Test GetNetworkTime
	t.Run("GetNetworkTime", func(t *testing.T) {
		timeInfo, err := GetNetworkTime(ids[0])
		if err != nil {
			// This might fail if the modem is not registered to a network
			t.Logf("Note: Failed to get network time: %v", err)
			t.Skip("Skipping test as network time is not available")
		}
		if timeInfo == nil {
			t.Error("Expected time info, got nil")
		} else {
			// Verify that the time info has the expected fields
			if timeInfo.NetworkTime == "" {
				t.Error("Expected network time, got empty string")
			}
			if timeInfo.Local == "" {
				t.Error("Expected local time, got empty string")
			}
		}
	})

	// Test GetNetworkTimeAsTime
	t.Run("GetNetworkTimeAsTime", func(t *testing.T) {
		networkTime, err := GetNetworkTimeAsTime(ids[0])
		if err != nil {
			// This might fail if the modem is not registered to a network
			t.Logf("Note: Failed to get network time as time.Time: %v", err)
			t.Skip("Skipping test as network time is not available")
		}

		// Verify that the time is not zero
		if networkTime.IsZero() {
			t.Error("Expected non-zero time, got zero time")
		}

		// Verify that the time is somewhat recent (within the last day)
		if time.Since(networkTime) > 24*time.Hour {
			t.Errorf("Expected recent time, got time more than 24 hours ago: %v", networkTime)
		}
	})
}
