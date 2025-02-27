package mmcli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// TimeInfo represents the network time information
type TimeInfo struct {
	NetworkTime string `json:"network-time"`
	Local       string `json:"local"`
}

// GetNetworkTime returns the current network time
func GetNetworkTime(modemID string) (*TimeInfo, error) {
	out, err := exec.Command("mmcli", "-m", modemID, "--time", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get network time: %w", err)
	}

	var response struct {
		Modem struct {
			Time TimeInfo `json:"time"`
		} `json:"modem"`
	}

	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse network time: %w", err)
	}

	return &response.Modem.Time, nil
}

// GetNetworkTimeAsTime returns the current network time as a time.Time object
func GetNetworkTimeAsTime(modemID string) (time.Time, error) {
	timeInfo, err := GetNetworkTime(modemID)
	if err != nil {
		return time.Time{}, err
	}

	// Parse the network time string
	// The format is expected to be ISO 8601, e.g. "2023-02-27T13:45:30+01:00"
	t, err := time.Parse(time.RFC3339, timeInfo.NetworkTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse network time string: %w", err)
	}

	return t, nil
}
