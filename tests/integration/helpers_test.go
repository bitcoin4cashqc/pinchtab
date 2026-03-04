package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:9867"

// requireOrchestrator skips the test if the orchestrator is not reachable.
func requireOrchestrator(t *testing.T) {
	t.Helper()
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/health", baseURL))
	if err != nil {
		t.Skipf("Orchestrator not reachable at %s (skipping): %v", baseURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
}
