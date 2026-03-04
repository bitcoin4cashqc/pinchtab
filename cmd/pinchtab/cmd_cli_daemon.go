package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pinchtab/pinchtab/internal/config"
)

// checkServerRunning checks if the pinchtab server is accessible
func checkServerRunning(base string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(base + "/")
	if err != nil {
		return fmt.Errorf("server not running on %s (use 'pinchtab start')", base)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("server not healthy (status %d)", resp.StatusCode)
	}
	return nil
}

// getServerURL returns the configured server URL
func getServerURL(cfg *config.RuntimeConfig) string {
	if envURL := os.Getenv("PINCHTAB_URL"); envURL != "" {
		return strings.TrimRight(envURL, "/")
	}
	return fmt.Sprintf("http://%s:%s", cfg.Bind, cfg.Port)
}

// getPIDFilePath returns the path to the PID file
func getPIDFilePath(cfg *config.RuntimeConfig) string {
	return filepath.Join(cfg.StateDir, "pinchtab.pid")
}

// cliStart starts the pinchtab server in the background
func cliStart(cfg *config.RuntimeConfig) {
	pidFile := getPIDFilePath(cfg)

	// Check if already running
	if pid, err := os.ReadFile(pidFile); err == nil {
		pidStr := strings.TrimSpace(string(pid))
		if pidNum, err := strconv.Atoi(pidStr); err == nil {
			// Check if process exists
			if proc, err := os.FindProcess(pidNum); err == nil {
				if err := proc.Signal(os.Signal(nil)); err == nil {
					fmt.Printf("✅ pinchtab already running (PID %d)\n", pidNum)
					fmt.Printf("   Dashboard: %s\n", getServerURL(cfg))
					return
				}
			}
		}
	}

	// Start server in background
	binary, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(binary)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to start server: %v\n", err)
		os.Exit(1)
	}

	// Write PID file
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", cmd.Process.Pid)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Started but failed to write PID file: %v\n", err)
	}

	// Wait a moment and check if server is up
	time.Sleep(2 * time.Second)
	baseURL := getServerURL(cfg)
	if err := checkServerRunning(baseURL); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Server started but not responding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ pinchtab started (PID %d)\n", cmd.Process.Pid)
	fmt.Printf("   Dashboard: %s\n", baseURL)
}

// cliStop stops the running pinchtab server
func cliStopDaemon(cfg *config.RuntimeConfig) {
	pidFile := getPIDFilePath(cfg)

	pid, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ No PID file found at %s\n", pidFile)
		fmt.Fprintf(os.Stderr, "   Server may not be running or was not started with 'pinchtab start'\n")
		os.Exit(1)
	}

	pidStr := strings.TrimSpace(string(pid))
	pidNum, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid PID in file: %s\n", pidStr)
		os.Exit(1)
	}

	proc, err := os.FindProcess(pidNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Process %d not found\n", pidNum)
		os.Remove(pidFile)
		os.Exit(1)
	}

	// Try to stop gracefully
	if err := proc.Signal(os.Interrupt); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to stop process %d: %v\n", pidNum, err)
		os.Exit(1)
	}

	// Wait for process to exit
	time.Sleep(1 * time.Second)

	// Verify it stopped
	if err := proc.Signal(os.Signal(nil)); err == nil {
		// Still running, force kill
		proc.Kill()
		time.Sleep(500 * time.Millisecond)
	}

	os.Remove(pidFile)
	fmt.Printf("✅ pinchtab stopped (PID %d)\n", pidNum)
}
