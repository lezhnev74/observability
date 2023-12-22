package observability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
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

// PollFpmPoolStatus connects to FPM directly and reads out the stats data
// It does not involve the web server.
func PollFpmPoolStatus(address, statusPath string) (*FpmStatus, error) {

	c := exec.Command("cgi-fcgi", "-bind", "-connect", address)
	c.Env = []string{
		fmt.Sprintf("SCRIPT_NAME=%s", statusPath),
		fmt.Sprintf("SCRIPT_FILENAME=%s", statusPath),
		"QUERY_STRING=json",
		"REQUEST_METHOD=GET",
	}

	cmdOut := bytes.NewBuffer(nil)
	c.Stdout = cmdOut

	err := c.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to read FPM stats: %w", err)
	}

	_, jsonStats, ok := bytes.Cut(cmdOut.Bytes(), []byte("\r\n\r\n"))
	if !ok {
		return nil, fmt.Errorf("unexpected FPM stats response: \n%s", cmdOut.Bytes())
	}

	var s FpmStatus
	err = json.Unmarshal(jsonStats, &s)
	if err != nil {
		return nil, fmt.Errorf("unable to parse FPM stats: %w, source string:\n%s", err, jsonStats)
	}

	return &s, nil
}
