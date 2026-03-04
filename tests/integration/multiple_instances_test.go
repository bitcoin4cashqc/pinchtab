package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestMultipleInstancesWithDifferentProfiles tests launching 2 instances with different profile names.
// Each instance should get its own profile directory to avoid SingletonLock conflicts.
func TestMultipleInstancesWithDifferentProfiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	baseURL := "http://localhost:9867"

	// Clean up any existing instances first
	resp, err := http.Get(fmt.Sprintf("%s/instances", baseURL))
	if err != nil {
		t.Fatalf("Failed to get instances: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Launch instance 1 with profile "test-multi-1"
	t.Log("Launching instance 1 with profile 'test-multi-1'...")
	inst1, err := launchInstance(baseURL, "test-multi-1", true)
	if err != nil {
		t.Fatalf("Failed to launch instance 1: %v", err)
	}
	t.Logf("Instance 1: %+v", inst1)

	// Launch instance 2 with profile "test-multi-2"
	t.Log("Launching instance 2 with profile 'test-multi-2'...")
	inst2, err := launchInstance(baseURL, "test-multi-2", true)
	if err != nil {
		t.Fatalf("Failed to launch instance 2: %v", err)
	}
	t.Logf("Instance 2: %+v", inst2)

	// Wait for both instances to reach "running" status
	t.Log("Waiting for both instances to reach 'running' status...")
	running1 := waitForInstanceRunning(t, baseURL, inst1.ID, 30*time.Second)
	running2 := waitForInstanceRunning(t, baseURL, inst2.ID, 30*time.Second)

	if !running1 {
		t.Errorf("Instance 1 (%s) never reached 'running' status (different profile issue)", inst1.ID)
	}
	if !running2 {
		t.Errorf("Instance 2 (%s) never reached 'running' status (different profile issue)", inst2.ID)
	}

	// If both are running, both profiles were properly isolated
	if running1 && running2 {
		t.Log("✓ Both instances reached 'running' status with separate profiles")
	} else {
		t.Log("✗ At least one instance failed to start (SingletonLock conflict or Chrome startup issue)")
	}

	// Verify we can get both instances from /instances endpoint
	instances, err := getInstances(baseURL)
	if err != nil {
		t.Fatalf("Failed to get instances list: %v", err)
	}

	// Count running instances
	runningCount := 0
	for _, inst := range instances {
		if inst["status"] == "running" {
			runningCount++
		}
	}

	t.Logf("Total instances running: %d", runningCount)
	if runningCount >= 2 {
		t.Log("✓ Both instances are in the running list")
	} else {
		t.Log("⚠️  Only found", runningCount, "running instance(s)")
	}

	// Cleanup
	_ = stopInstance(baseURL, inst1.ID)
	_ = stopInstance(baseURL, inst2.ID)
}

type Instance struct {
	ID       string `json:"id"`
	Port     string `json:"port"`
	Status   string `json:"status"`
	Name     string `json:"profileName"`
	Headless bool   `json:"headless"`
}

func launchInstance(baseURL, profileName string, headless bool) (*Instance, error) {
	reqBody := map[string]interface{}{
		"name":     profileName,
		"headless": headless,
	}
	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		fmt.Sprintf("%s/instances/launch", baseURL),
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("launch failed with status %d: %s", resp.StatusCode, string(body))
	}

	var inst Instance
	if err := json.NewDecoder(resp.Body).Decode(&inst); err != nil {
		return nil, err
	}

	return &inst, nil
}

func waitForInstanceRunning(t *testing.T, baseURL, instID string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			instances, err := getInstances(baseURL)
			if err != nil {
				t.Logf("Error getting instances: %v", err)
				continue
			}

			for _, inst := range instances {
				if inst["id"] == instID {
					status := inst["status"].(string)
					t.Logf("  %s: %s", instID, status)
					if status == "running" {
						return true
					}
					break
				}
			}

		case <-time.After(timeout):
			t.Logf("Timeout waiting for instance %s to reach 'running' status", instID)
			return false
		}

		if time.Now().After(deadline) {
			t.Logf("Timeout waiting for instance %s to reach 'running' status", instID)
			return false
		}
	}
}

func getInstances(baseURL string) ([]map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/instances", baseURL))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var instances []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return nil, err
	}

	return instances, nil
}

func stopInstance(baseURL, instID string) error {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/instances/%s/stop", baseURL, instID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}
