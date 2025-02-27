package mmcli

import (
	"fmt"
	"os/exec"
	"strings"
)

// ConnectSettings represents the settings for a simple connection
type ConnectSettings struct {
	APN          string // Access Point Name
	User         string // Username for authentication
	Password     string // Password for authentication
	IPType       string // IP type (ipv4, ipv6, ipv4v6)
	Number       string // Number to dial (for PPP connections)
	AllowRoaming bool   // Whether to allow roaming
	PIN          string // SIM PIN if required
}

// Connect establishes a connection with the specified settings
func Connect(modemID string, settings ConnectSettings) error {
	// Build the settings string
	var settingsParams []string
	if settings.APN != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("apn=%s", settings.APN))
	}
	if settings.User != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("user=%s", settings.User))
	}
	if settings.Password != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("password=%s", settings.Password))
	}
	if settings.IPType != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("ip-type=%s", settings.IPType))
	}
	if settings.Number != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("number=%s", settings.Number))
	}
	if settings.AllowRoaming {
		settingsParams = append(settingsParams, "allow-roaming=yes")
	}
	if settings.PIN != "" {
		settingsParams = append(settingsParams, fmt.Sprintf("pin=%s", settings.PIN))
	}

	// Construct the command
	args := []string{"-m", modemID}

	// Add the settings string to the arguments if we have any settings
	if len(settingsParams) > 0 {
		settingsStr := strings.Join(settingsParams, ",")
		args = append(args, fmt.Sprintf("--simple-connect=\"%s\"", settingsStr))
	} else {
		args = append(args, "--simple-connect")
	}

	// Execute the command
	_, err := exec.Command("mmcli", args...).Output()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	return nil
}

// Disconnect disconnects all connected bearers
func Disconnect(modemID string) error {
	_, err := exec.Command("mmcli", "-m", modemID, "--simple-disconnect").Output()
	if err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}

	return nil
}

// ConnectWithAPN is a convenience function to connect with just an APN
func ConnectWithAPN(modemID string, apn string) error {
	settings := ConnectSettings{
		APN: apn,
	}
	return Connect(modemID, settings)
}

// ConnectWithAuth is a convenience function to connect with APN, username and password
func ConnectWithAuth(modemID string, apn, user, password string) error {
	settings := ConnectSettings{
		APN:      apn,
		User:     user,
		Password: password,
	}
	return Connect(modemID, settings)
}
