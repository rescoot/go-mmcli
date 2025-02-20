package mmcli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ModemList struct {
	ModemList []string `json:"modem-list"`
}

type ModemManager struct {
	Modem Modem `json:"modem"`
}

type Modem struct {
	ThreeGPP ThreeGPP         `json:"3gpp"`
	CDMA     CDMA             `json:"cdma"`
	DBusPath string           `json:"dbus-path"`
	Generic  ModemGenericInfo `json:"generic"`
}

type ThreeGPP struct {
	EnabledLocks      []string `json:"enabled-locks"`
	EPS               EPSInfo  `json:"eps"`
	IMEI              string   `json:"imei"`
	OperatorCode      string   `json:"operator-code"`
	OperatorName      string   `json:"operator-name"`
	PCO               string   `json:"pco"`
	RegistrationState string   `json:"registration-state"`
}

// EPSInfo contains EPS (Evolved Packet System) information
type EPSInfo struct {
	InitialBearer   EPSBearer `json:"initial-bearer"`
	UEModeOperation string    `json:"ue-mode-operation"`
}

// EPSBearer contains bearer settings
type EPSBearer struct {
	DBusPath string         `json:"dbus-path"`
	Settings BearerSettings `json:"settings"`
}

type BearerSettings struct {
	APN      string `json:"apn"`
	IPType   string `json:"ip-type"`
	Password string `json:"password"`
	User     string `json:"user"`
}

type CDMA struct {
	ActivationState         string `json:"activation-state"`
	CDMA1xRegistrationState string `json:"cdma1x-registration-state"`
	ESN                     string `json:"esn"`
	EVDORegistrationState   string `json:"evdo-registration-state"`
	MEID                    string `json:"meid"`
	NID                     string `json:"nid"`
	SID                     string `json:"sid"`
}

type ModemGenericInfo struct {
	AccessTechnologies    []string      `json:"access-technologies"`
	Bearers               []string      `json:"bearers"`
	CarrierConfiguration  string        `json:"carrier-configuration"`
	CurrentBands          []string      `json:"current-bands"`
	CurrentCapabilities   []string      `json:"current-capabilities"`
	CurrentModes          string        `json:"current-modes"`
	Device                string        `json:"device"`
	DeviceIdentifier      string        `json:"device-identifier"`
	Drivers               []string      `json:"drivers"`
	EquipmentIdentifier   string        `json:"equipment-identifier"`
	HardwareRevision      string        `json:"hardware-revision"`
	Manufacturer          string        `json:"manufacturer"`
	Model                 string        `json:"model"`
	OwnNumbers            []string      `json:"own-numbers"`
	Plugin                string        `json:"plugin"`
	Ports                 []string      `json:"ports"`
	PowerState            string        `json:"power-state"`
	PrimaryPort           string        `json:"primary-port"`
	Revision              string        `json:"revision"`
	SignalQuality         SignalQuality `json:"signal-quality"`
	SIM                   string        `json:"sim"`
	State                 string        `json:"state"`
	StateFailedReason     string        `json:"state-failed-reason"`
	SupportedBands        []string      `json:"supported-bands"`
	SupportedCapabilities []string      `json:"supported-capabilities"`
	SupportedIPFamilies   []string      `json:"supported-ip-families"`
	SupportedModes        []string      `json:"supported-modes"`
	UnlockRequired        string        `json:"unlock-required"`
	UnlockRetries         []string      `json:"unlock-retries"`
}

type SignalQuality struct {
	Recent string `json:"recent"`
	Value  string `json:"value"`
}

type SIMInfo struct {
	DBusPath   string        `json:"dbus-path"`
	Properties SIMProperties `json:"properties"`
}

type SIMProperties struct {
	Active           string   `json:"active"`
	EID              string   `json:"eid"`
	EmergencyNumbers []string `json:"emergency-numbers"`
	ICCID            string   `json:"iccid"`
	IMSI             string   `json:"imsi"`
	OperatorCode     string   `json:"operator-code"`
	OperatorName     string   `json:"operator-name"`
}

// ListModems returns a list of all available modems with their IDs
func ListModems() ([]string, error) {
	out, err := exec.Command("mmcli", "-J", "-L").Output()
	if err != nil {
		return nil, fmt.Errorf("mmcli list error: %w", err)
	}

	var list ModemList
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("failed to parse modem list: %w", err)
	}

	return list.ModemList, nil
}

// GetModemIDs returns a list of numeric modem IDs
func GetModemIDs() ([]string, error) {
	modems, err := ListModems()
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(modems))
	for i, path := range modems {
		parts := strings.Split(path, "/")
		if len(parts) < 6 {
			return nil, fmt.Errorf("invalid modem path format: %s", path)
		}
		ids[i] = parts[5]
	}

	return ids, nil
}

// GetFirstModemID returns the ID of the first available modem
func GetFirstModemID() (string, error) {
	ids, err := GetModemIDs()
	if err != nil {
		return "", err
	}

	if len(ids) == 0 {
		return "", fmt.Errorf("no modem found")
	}

	return ids[0], nil
}

// GetModemDetails returns details for a specific modem by ID
func GetModemDetails(modemID string) (*ModemManager, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get modem details: %w", err)
	}

	return Parse(out)
}

func ResetModem(modemID string) (bool, error) {
	_, err := exec.Command("mmcli", "-m", modemID, "--reset").Output()
	if err != nil {
		return false, fmt.Errorf("failed to reset modem: %w", err)
	}

	return true, nil
}

func GetSIMInfo(modemID string) (*SIMInfo, error) {
	out, err := exec.Command("mmcli", "-i", modemID, "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get SIM info: %w", err)
	}

	var simInfo SIMInfo
	if err := json.Unmarshal(out, &simInfo); err != nil {
		return nil, fmt.Errorf("failed to parse SIM info: %w", err)
	}

	return &simInfo, nil
}

// Parse parses mmcli JSON output into a ModemManager struct
func Parse(data []byte) (*ModemManager, error) {
	var mm ModemManager
	if err := json.Unmarshal(data, &mm); err != nil {
		return nil, fmt.Errorf("failed to parse mmcli output: %w", err)
	}
	return &mm, nil
}

// Helper methods for ModemManager

// IsConnected returns true if the modem is in connected state
func (mm *ModemManager) IsConnected() bool {
	return mm.Modem.Generic.State == "connected"
}

// SignalStrength returns the signal strength as an integer percentage
func (mm *ModemManager) SignalStrength() (int, error) {
	var strength int
	_, err := fmt.Sscanf(mm.Modem.Generic.SignalQuality.Value, "%d", &strength)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signal strength: %w", err)
	}
	return strength, nil
}

// GetPortByType returns the port name for a given type (qmi, at, net, etc)
func (mm *ModemManager) GetPortByType(portType string) string {
	portType = strings.ToLower(portType)
	for _, port := range mm.Modem.Generic.Ports {
		if strings.Contains(strings.ToLower(port), "("+portType+")") {
			parts := strings.Split(port, " ")
			if len(parts) > 0 {
				return parts[0]
			}
		}
	}
	return ""
}

// IsSimLocked returns true if the modem requires a SIM PIN
func (mm *ModemManager) IsSimLocked() bool {
	return strings.HasPrefix(mm.Modem.Generic.UnlockRequired, "sim-")
}

// GetCurrentAccessTechnology returns the current access technology (2G/3G/4G)
func (mm *ModemManager) GetCurrentAccessTechnology() string {
	for _, tech := range mm.Modem.Generic.AccessTechnologies {
		switch strings.ToLower(tech) {
		case "lte":
			return "4G"
		case "umts", "hspa", "hspa+":
			return "3G"
		case "gsm", "gprs", "edge":
			return "2G"
		}
	}
	return "Unknown"
}

// GetOperatorInfo returns the operator name and code
func (mm *ModemManager) GetOperatorInfo() (name string, code string) {
	return mm.Modem.ThreeGPP.OperatorName, mm.Modem.ThreeGPP.OperatorCode
}

// RemainingUnlockRetries returns the remaining unlock attempts for a given lock type
func (mm *ModemManager) RemainingUnlockRetries(lockType string) int {
	for _, retry := range mm.Modem.Generic.UnlockRetries {
		if strings.HasPrefix(retry, lockType) {
			var attempts int
			fmt.Sscanf(retry, "%s (%d)", &lockType, &attempts)
			return attempts
		}
	}
	return 0
}

// IsIPv6Supported returns true if the modem supports IPv6
func (mm *ModemManager) IsIPv6Supported() bool {
	for _, family := range mm.Modem.Generic.SupportedIPFamilies {
		if family == "ipv6" || family == "ipv4v6" {
			return true
		}
	}
	return false
}

// GetAllPorts returns a map of port types to port names
func (mm *ModemManager) GetAllPorts() map[string]string {
	ports := make(map[string]string)
	for _, port := range mm.Modem.Generic.Ports {
		parts := strings.Split(port, " ")
		if len(parts) >= 2 {
			portType := strings.Trim(parts[1], "()")
			ports[portType] = parts[0]
		}
	}
	return ports
}
