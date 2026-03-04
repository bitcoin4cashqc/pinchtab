package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func cliHealth(client *http.Client, base, token string) {
	result := doGet(client, base, token, "/health", nil)
	if status, ok := result["status"].(string); ok && status == "ok" {
		fmt.Println("✅ Server is healthy")
		if strategy, ok := result["strategy"].(string); ok {
			fmt.Printf("   Strategy: %s\n", strategy)
		}
	} else {
		fmt.Println("❌ Server health check failed")
		os.Exit(1)
	}
}

func cliProfiles(client *http.Client, base, token string) {
	result := doGet(client, base, token, "/profiles", nil)

	if profiles, ok := result["profiles"].([]interface{}); ok && len(profiles) > 0 {
		fmt.Println("\n📋 Available Profiles:")
		fmt.Println()
		for _, prof := range profiles {
			if m, ok := prof.(map[string]any); ok {
				name, _ := m["name"].(string)
				fmt.Printf("  👤 %s\n", name)
			}
		}
		fmt.Println()
	} else {
		fmt.Println("No profiles available")
	}
}

func cliInstances(client *http.Client, base, token string) {
	// /instances returns an array directly - just print it as JSON
	body := doGetRaw(client, base, token, "/instances", nil)
	fmt.Println(string(body))
}

func cliLaunch(client *http.Client, base, token string, args []string) {
	body := map[string]any{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--headed":
			body["mode"] = "headed"
		case "--port":
			if i+1 < len(args) {
				i++
				body["port"] = args[i]
			}
		default:
			// Positional arg = profile name
			if !strings.HasPrefix(args[i], "-") {
				body["name"] = args[i]
			}
		}
	}

	result := doPost(client, base, token, "/instances/launch", body)
	id, _ := result["id"].(string)
	port, _ := result["port"].(string)
	fmt.Printf("🚀 Launched %s on port %s\n", id, port)
}

func cliStop(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: pinchtab stop <instance-id>")
		os.Exit(1)
	}
	doPost(client, base, token, fmt.Sprintf("/instances/%s/stop", args[0]), nil)
	fmt.Printf("⏹️  Stopped %s\n", args[0])
}

func cliTabs(client *http.Client, base, token string) {
	// /tabs returns an array directly - just print it as JSON
	body := doGetRaw(client, base, token, "/tabs", nil)
	fmt.Println(string(body))
}

func cliOpen(client *http.Client, base, token string, args []string) {
	tabURL := ""
	if len(args) > 0 {
		tabURL = args[0]
	}

	body := map[string]any{"action": "new"}
	if tabURL != "" {
		body["url"] = tabURL
	}

	result := doPost(client, base, token, "/tab", body)
	id, _ := result["tabId"].(string)
	resultURL, _ := result["url"].(string)
	if resultURL == "" {
		resultURL = "about:blank"
	}
	fmt.Printf("📑 Opened [%s] → %s\n", id, resultURL)
}

func cliClose(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: pinchtab close <tab-id>")
		os.Exit(1)
	}
	doDelete(client, base, token, fmt.Sprintf("/tabs/%s", args[0]))
	fmt.Printf("🗑️  Closed %s\n", args[0])
}
