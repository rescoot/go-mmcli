package mmcli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// MessagingStatus represents the status of messaging support
type MessagingStatus struct {
	DefaultStorages   []string `json:"default-storages"`
	SupportedStorages []string `json:"supported-storages"`
}

// SMSList represents a list of SMS messages
type SMSList struct {
	SMS []string `json:"modem.messaging.sms"`
}

// SMSInfo represents information about an SMS message
type SMSInfo struct {
	DBusPath   string        `json:"dbus-path"`
	Properties SMSProperties `json:"properties"`
}

// SMSProperties contains the properties of an SMS message
type SMSProperties struct {
	Class              int    `json:"class"`
	DeliveryReportReq  bool   `json:"delivery-report-request"`
	DeliveryState      string `json:"delivery-state"`
	DischargeTimestamp string `json:"discharge-timestamp"`
	Number             string `json:"number"`
	PDU                string `json:"pdu"`
	SMSC               string `json:"smsc"`
	State              string `json:"state"`
	Storage            string `json:"storage"`
	Teleservice        string `json:"teleservice-id"`
	Text               string `json:"text"`
	Timestamp          string `json:"timestamp"`
	Validity           string `json:"validity"`
	Data               []byte `json:"data"`
}

// SMSCreateSettings represents the settings for creating a new SMS
type SMSCreateSettings struct {
	Number            string // Destination phone number
	Text              string // Message text
	SMSC              string // SMS service center number (optional)
	Validity          string // Validity period (optional)
	Class             int    // Message class (optional)
	DeliveryReportReq bool   // Request delivery report (optional)
}

// GetMessagingStatus returns the status of messaging support
func GetMessagingStatus(modemID string) (*MessagingStatus, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "--messaging-status", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get messaging status: %w", err)
	}

	var response struct {
		Modem struct {
			Messaging MessagingStatus `json:"messaging"`
		} `json:"modem"`
	}

	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse messaging status: %w", err)
	}

	return &response.Modem.Messaging, nil
}

// ListSMS returns a list of SMS messages
func ListSMS(modemID string) ([]string, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "--messaging-list-sms", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list SMS messages: %w", err)
	}

	var list SMSList
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("failed to parse SMS list: %w", err)
	}

	return list.SMS, nil
}

// GetSMSInfo returns information about a specific SMS message
func GetSMSInfo(smsID string) (*SMSInfo, error) {
	out, err := exec.Command("mmcli", "-s", smsID, "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get SMS info: %w", err)
	}

	var response struct {
		SMS SMSInfo `json:"sms"`
	}

	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse SMS info: %w", err)
	}

	return &response.SMS, nil
}

// CreateSMS creates a new SMS message
func CreateSMS(modemID string, settings SMSCreateSettings) (string, error) {
	// Build the settings string
	var settingsParams []string
	if settings.Number != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("number=%s", settings.Number))
	}
	if settings.Text != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("text=%s", settings.Text))
	}
	if settings.SMSC != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("smsc=%s", settings.SMSC))
	}
	if settings.Validity != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("validity=%s", settings.Validity))
	}
	if settings.Class != 0 {
		settingsParams = append(settingsParams, fmt.Sprintf("class=%d", settings.Class))
	}
	if settings.DeliveryReportReq {
		settingsParams = append(settingsParams, "delivery-report-request=yes")
	}

	// Construct the command
	args := []string{"-m", modemID}

	// Add the settings string to the arguments
	settingsStr := strings.Join(settingsParams, ",")
	args = append(args, fmt.Sprintf("--messaging-create-sms=\"%s\"", settingsStr))

	// Execute the command
	out, err := exec.Command("mmcli", args...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to create SMS: %w", err)
	}

	// Parse the output to get the SMS ID
	// Expected format: "Successfully created new SMS: /org/freedesktop/ModemManager1/SMS/X"
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Successfully created new SMS:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				smsPath := strings.TrimSpace(parts[1])
				// Extract the SMS ID from the path
				pathParts := strings.Split(smsPath, "/")
				if len(pathParts) > 0 {
					return pathParts[len(pathParts)-1], nil
				}
			}
		}
	}

	return "", fmt.Errorf("failed to parse SMS ID from output")
}

// SendSMS sends an SMS message
func SendSMS(smsID string) error {
	_, err := exec.Command("mmcli", "-s", smsID, "--send").Output()
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

// StoreSMS stores an SMS message in the device
func StoreSMS(smsID string) error {
	_, err := exec.Command("mmcli", "-s", smsID, "--store").Output()
	if err != nil {
		return fmt.Errorf("failed to store SMS: %w", err)
	}

	return nil
}

// StoreSMSInStorage stores an SMS message in the specified storage
func StoreSMSInStorage(smsID string, storage string) error {
	_, err := exec.Command("mmcli", "-s", smsID, fmt.Sprintf("--store-in-storage=%s", storage)).Output()
	if err != nil {
		return fmt.Errorf("failed to store SMS in storage %s: %w", storage, err)
	}

	return nil
}

// DeleteSMS deletes an SMS message
func DeleteSMS(modemID string, smsID string) error {
	_, err := exec.Command("mmcli", "-m", modemID, fmt.Sprintf("--messaging-delete-sms=%s", smsID)).Output()
	if err != nil {
		return fmt.Errorf("failed to delete SMS: %w", err)
	}

	return nil
}

// CreateAndSendSMS creates and sends an SMS message in one step
func CreateAndSendSMS(modemID string, number string, text string) error {
	settings := SMSCreateSettings{
		Number: number,
		Text:   text,
	}

	smsID, err := CreateSMS(modemID, settings)
	if err != nil {
		return err
	}

	return SendSMS(smsID)
}
