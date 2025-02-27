package mmcli

import (
	"os/exec"
	"testing"
)

func TestSMSFunctions(t *testing.T) {
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

	// Test GetMessagingStatus
	t.Run("GetMessagingStatus", func(t *testing.T) {
		status, err := GetMessagingStatus(ids[0])
		if err != nil {
			t.Errorf("Failed to get messaging status: %v", err)
		}
		if status == nil {
			t.Error("Expected messaging status, got nil")
		} else {
			// Verify that the status has the expected fields
			if status.DefaultStorages == nil {
				t.Error("Expected default storages, got nil")
			}
			if status.SupportedStorages == nil {
				t.Error("Expected supported storages, got nil")
			}
		}
	})

	// Test ListSMS
	t.Run("ListSMS", func(t *testing.T) {
		smslist, err := ListSMS(ids[0])
		if err != nil {
			t.Errorf("Failed to list SMS messages: %v", err)
		}
		if smslist == nil {
			t.Error("Expected SMS list, got nil")
		}
		// Note: We can't verify the actual values as they depend on the modem's state
	})

	// Test SMS creation and sending
	// Note: We're not actually creating/sending SMS here to avoid side effects
	t.Run("CreateSMSSettings", func(t *testing.T) {
		// Just verify that the settings struct works correctly
		settings := SMSCreateSettings{
			Number:            "+1234567890",
			Text:              "Test message",
			SMSC:              "+9876543210",
			Validity:          "1h",
			Class:             1,
			DeliveryReportReq: true,
		}

		// Verify the settings are set correctly
		if settings.Number != "+1234567890" {
			t.Errorf("Expected Number '+1234567890', got '%s'", settings.Number)
		}
		if settings.Text != "Test message" {
			t.Errorf("Expected Text 'Test message', got '%s'", settings.Text)
		}
		if settings.SMSC != "+9876543210" {
			t.Errorf("Expected SMSC '+9876543210', got '%s'", settings.SMSC)
		}
		if settings.Validity != "1h" {
			t.Errorf("Expected Validity '1h', got '%s'", settings.Validity)
		}
		if settings.Class != 1 {
			t.Errorf("Expected Class 1, got %d", settings.Class)
		}
		if !settings.DeliveryReportReq {
			t.Errorf("Expected DeliveryReportReq true, got false")
		}
	})
}
