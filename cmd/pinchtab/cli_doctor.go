package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func cliDoctor(client *http.Client, base string, token string, args []string) {
	fmt.Println("🦀 Pinchtab Doctor - Setup & Diagnostics")
	fmt.Println("")

	passed := 0
	failed := 0

	// 1. Check git hooks
	fmt.Print("Checking git hooks configuration... ")
	if checkAndSetupGitHooks() {
		fmt.Println("✅")
		passed++
	} else {
		fmt.Println("⚠️")
		failed++
	}

	// 2. Check server connection
	fmt.Print("Checking server connection... ")
	result := doGet(client, base, token, "/health", nil)
	if result == nil {
		fmt.Println("✅")
		passed++
	} else {
		fmt.Println("❌")
		failed++
	}

	// 3. Check Go installation
	fmt.Print("Checking Go installation... ")
	if checkGo() {
		fmt.Println("✅")
		passed++
	} else {
		fmt.Println("❌")
		failed++
	}

	// 4. Check Chrome/Chromium
	fmt.Print("Checking Chrome/Chromium... ")
	if checkChrome() {
		fmt.Println("✅")
		passed++
	} else {
		fmt.Println("⚠️ (optional, but required for full functionality)")
		failed++
	}

	fmt.Println("")
	fmt.Printf("Results: %d passed, %d issues\n", passed, failed)

	if failed == 0 {
		fmt.Println("")
		fmt.Println("✅ All checks passed! Pinchtab is ready to use.")
	}
}

// checkAndSetupGitHooks configures git hooks for the repository
func checkAndSetupGitHooks() bool {
	// Check if .githooks directory exists
	if _, err := os.Stat(".githooks"); err != nil {
		fmt.Println("(not in repo root, skipping)")
		return true
	}

	// Check current hooks path
	cmd := exec.Command("git", "config", "core.hooksPath")
	output, err := cmd.Output()
	if err == nil && string(output) == ".githooks\n" {
		return true // Already configured
	}

	// Configure git hooks
	cmd = exec.Command("git", "config", "core.hooksPath", ".githooks")
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

// checkGo verifies Go is installed
func checkGo() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

// checkChrome verifies Chrome/Chromium is available
func checkChrome() bool {
	// Try common Chrome/Chromium locations
	paths := []string{
		"/usr/bin/google-chrome",                                       // Linux
		"/snap/bin/chromium",                                           // Linux snap
		"/opt/google/chrome/google-chrome",                             // Linux custom
		"/usr/bin/chromium-browser",                                    // Linux
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", // macOS
		"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",   // Windows
		"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}
