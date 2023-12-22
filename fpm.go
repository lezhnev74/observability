package observability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FpmStatus is a mapping from what Fpm status page returns
// see https://www.php.net/manual/en/fpm.status.php
type FpmStatus struct {
	// The name of the FPM process pool.
	Name string `json:"pool"`
	// The total number of accepted connections.
	TotalAccepted int `json:"accepted conn"`
	// The number of requests (backlog) currently waiting for a free process.
	CurrentWaitingClients int `json:"listen queue"`
	// The maximum number of requests seen in the listen queue at any one time.
	MaxWaitingClients int `json:"max listen queue"`
	// The maximum number of concurrently active processes.
	MaxActiveProcesses int `json:"max active processes"`
	// he number of processes that are currently idle (waiting for requests).
	IdleProcesses int `json:"idle processes"`
	// The number of processes that are currently processing requests.
	ActiveProcesses int `json:"active processes"`
	// The total number of requests that have hit the configured request_slowlog_timeout.
	SlowRequests int `json:"slow requests"`
}

// PollFpmStatus connects to FPM directly and reads out the stats data
// It does not involve the web server.
func PollFpmStatus(FpmPoolStatsUrl string) (*FpmStatus, error) {
	jsonUrl := fmt.Sprintf("%s?json", FpmPoolStatsUrl)
	resp, err := http.Get(jsonUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to get FPM stats: %s", err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body of FPM stats: %s", err)
	}

	var s FpmStatus
	err = json.Unmarshal(buf.Bytes(), &s)
	if err != nil {
		return nil, fmt.Errorf("unable to parse FPM stats: %s", err)
	}

	return &s, nil
}
