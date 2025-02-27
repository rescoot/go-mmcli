package mmcli

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// LocationStatus represents the status of location gathering
type LocationStatus struct {
	Capabilities []string `json:"capabilities"`
	Enabled      []string `json:"enabled"`
	Signals      string   `json:"signals"`
	GPS          GPSInfo  `json:"gps"`
}

// GPSInfo contains GPS-specific information
type GPSInfo struct {
	Assistance        []string `json:"assistance"`
	AssistanceServers []string `json:"assistance-servers"`
	RefreshRate       string   `json:"refresh-rate"`
	SuplServer        string   `json:"supl-server"`
}

// LocationInfo represents location information from various sources
type LocationInfo struct {
	ThreeGPP ThreeGPPLocation `json:"3gpp"`
	GPS      GPSLocation      `json:"gps"`
	CDMABS   CDMALocation     `json:"cdma-bs"`
}

// ThreeGPPLocation contains 3GPP-specific location information
type ThreeGPPLocation struct {
	MCC string `json:"mcc"`
	MNC string `json:"mnc"`
	LAC string `json:"lac"`
	CID string `json:"cid"`
	TAC string `json:"tac"`
}

// GPSLocation contains GPS-specific location information
type GPSLocation struct {
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	Altitude  string   `json:"altitude"`
	UTC       string   `json:"utc"`
	NMEA      []string `json:"nmea"`
}

// CDMALocation contains CDMA-specific location information
type CDMALocation struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// GetLocationStatus returns the current status of location gathering
func GetLocationStatus(modemID string) (*LocationStatus, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "--location-status", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get location status: %w", err)
	}

	var response struct {
		Modem struct {
			Location LocationStatus `json:"location"`
		} `json:"modem"`
	}

	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse location status: %w", err)
	}

	return &response.Modem.Location, nil
}

// GetLocation returns the current location information
func GetLocation(modemID string) (*LocationInfo, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "--location-get", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	var response struct {
		Modem struct {
			Location LocationInfo `json:"location"`
		} `json:"modem"`
	}

	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse location info: %w", err)
	}

	return &response.Modem.Location, nil
}

// EnableLocationGathering enables a specific location gathering method
func EnableLocationGathering(modemID string, method string) error {
	var option string
	switch method {
	case "3gpp":
		option = "--location-enable-3gpp"
	case "agps-msa":
		option = "--location-enable-agps-msa"
	case "agps-msb":
		option = "--location-enable-agps-msb"
	case "gps-nmea":
		option = "--location-enable-gps-nmea"
	case "gps-raw":
		option = "--location-enable-gps-raw"
	case "cdma-bs":
		option = "--location-enable-cdma-bs"
	case "gps-unmanaged":
		option = "--location-enable-gps-unmanaged"
	default:
		return fmt.Errorf("unsupported location method: %s", method)
	}

	_, err := exec.Command("mmcli", "-m", modemID, option).Output()
	if err != nil {
		return fmt.Errorf("failed to enable location gathering (%s): %w", method, err)
	}

	return nil
}

// DisableLocationGathering disables a specific location gathering method
func DisableLocationGathering(modemID string, method string) error {
	var option string
	switch method {
	case "3gpp":
		option = "--location-disable-3gpp"
	case "agps-msa":
		option = "--location-disable-agps-msa"
	case "agps-msb":
		option = "--location-disable-agps-msb"
	case "gps-nmea":
		option = "--location-disable-gps-nmea"
	case "gps-raw":
		option = "--location-disable-gps-raw"
	case "cdma-bs":
		option = "--location-disable-cdma-bs"
	case "gps-unmanaged":
		option = "--location-disable-gps-unmanaged"
	default:
		return fmt.Errorf("unsupported location method: %s", method)
	}

	_, err := exec.Command("mmcli", "-m", modemID, option).Output()
	if err != nil {
		return fmt.Errorf("failed to disable location gathering (%s): %w", method, err)
	}

	return nil
}

// SetSuplServer sets the SUPL server address for A-GPS
func SetSuplServer(modemID string, address string) error {
	_, err := exec.Command("mmcli", "-m", modemID, "--location-set-supl-server="+address).Output()
	if err != nil {
		return fmt.Errorf("failed to set SUPL server: %w", err)
	}

	return nil
}

// SetGpsRefreshRate sets the GPS refresh rate in seconds
func SetGpsRefreshRate(modemID string, rate int) error {
	_, err := exec.Command("mmcli", "-m", modemID, fmt.Sprintf("--location-set-gps-refresh-rate=%d", rate)).Output()
	if err != nil {
		return fmt.Errorf("failed to set GPS refresh rate: %w", err)
	}

	return nil
}

// EnableLocationSignals enables location update signaling in DBus property
func EnableLocationSignals(modemID string) error {
	_, err := exec.Command("mmcli", "-m", modemID, "--location-set-enable-signal").Output()
	if err != nil {
		return fmt.Errorf("failed to enable location signals: %w", err)
	}

	return nil
}

// DisableLocationSignals disables location update signaling in DBus property
func DisableLocationSignals(modemID string) error {
	_, err := exec.Command("mmcli", "-m", modemID, "--location-set-disable-signal").Output()
	if err != nil {
		return fmt.Errorf("failed to disable location signals: %w", err)
	}

	return nil
}
