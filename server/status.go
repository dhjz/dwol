package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

type MachineStatus struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Status string `json:"status"`
	Method string `json:"method"`
}

func checkPing(host string) bool {
	if host == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := fmt.Sprintf("ping -n 1 -w 2000 %s", host)

	type result struct {
		success bool
	}

	resultChan := make(chan result, 1)

	go func() {
		var cmdResult result
		if _, err := execCommand(ctx, cmd); err == nil {
			cmdResult.success = true
		}
		resultChan <- cmdResult
	}()

	select {
	case res := <-resultChan:
		return res.success
	case <-ctx.Done():
		return false
	}
}

func checkPort(host string, port int, timeout time.Duration) bool {
	if host == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func execCommand(ctx context.Context, cmd string) ([]byte, error) {
	type result struct {
		output []byte
		err    error
	}

	resultChan := make(chan result, 1)

	go func() {
		var res result
		res.output, res.err = exec.CommandContext(ctx, "cmd", "/c", cmd).CombinedOutput()
		resultChan <- res
	}()

	select {
	case res := <-resultChan:
		return res.output, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func checkMachineStatus(machine Machine) MachineStatus {
	status := MachineStatus{
		ID:   machine.ID,
		Host: machine.Host,
	}

	if machine.Host == "" {
		status.Status = "unknown"
		status.Method = "none"
		return status
	}

	type checkResult struct {
		status string
		method string
	}

	resultChan := make(chan checkResult, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if checkPing(machine.Host) {
			resultChan <- checkResult{status: "online", method: "ping"}
		}
	}()

	go func() {
		defer wg.Done()
		if checkPort(machine.Host, 22, 2*time.Second) {
			resultChan <- checkResult{status: "online", method: "ssh"}
		}
	}()

	go func() {
		defer wg.Done()
		if checkPort(machine.Host, 3389, 2*time.Second) {
			resultChan <- checkResult{status: "online", method: "rdp"}
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		status.Status = result.status
		status.Method = result.method
		return status
	}

	status.Status = "offline"
	status.Method = "none"
	return status
}

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	machines := GetMachines()
	results := make([]MachineStatus, len(machines))

	var wg sync.WaitGroup
	wg.Add(len(machines))

	for i, machine := range machines {
		go func(idx int, m Machine) {
			defer wg.Done()
			results[idx] = checkMachineStatus(m)
		}(i, machine)
	}

	wg.Wait()

	sendJSON(w, results)
}

func HandleTestPort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Host == "" || req.Port <= 0 || req.Port > 65535 {
		sendError(w, "Invalid host or port", http.StatusBadRequest)
		return
	}

	connected := checkPort(req.Host, req.Port, 3*time.Second)

	sendJSON(w, map[string]interface{}{
		"host":      req.Host,
		"port":      req.Port,
		"connected": connected,
	})
}
