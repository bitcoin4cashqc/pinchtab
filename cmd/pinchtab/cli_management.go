package main

import (
	"fmt"
	"net/http"
	"os"
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
	result := doGet(client, base, token, "/instances", nil)

	if instances, ok := result["instances"].([]interface{}); ok {
		if len(instances) == 0 {
			fmt.Println("No instances running")
			fmt.Println("\nTo launch an instance:")
			fmt.Println("  1. Start dashboard: pinchtab")
			fmt.Println("  2. Open browser: http://localhost:9867/dashboard")
			fmt.Println("  3. Click 'Profiles' → select profile → 'Launch'")
			return
		}

		fmt.Println("\n🚀 Running Instances:")
		fmt.Println()

		for _, inst := range instances {
			if m, ok := inst.(map[string]any); ok {
				id, _ := m["id"].(string)
				port, _ := m["port"].(string)
				status, _ := m["status"].(string)
				headless, _ := m["headless"].(bool)

				mode := "headless"
				if !headless {
					mode = "headed"
				}

				icon := "▶️"
				if status != "running" {
					icon = "⏸️"
				}

				fmt.Printf("  %s %s (port %s, %s)\n", icon, id, port, mode)
				if port != "" && status == "running" {
					fmt.Printf("     → http://localhost:%s\n", port)
				}
			}
		}
		fmt.Println()
	} else {
		fmt.Println("Failed to get instances")
		os.Exit(1)
	}
}

func cliTabs(client *http.Client, base, token string) {
	result := doGet(client, base, token, "/tabs", nil)

	if tabs, ok := result["tabs"].([]interface{}); ok {
		if len(tabs) == 0 {
			fmt.Println("No tabs open across all instances")
			return
		}

		fmt.Println("\n📑 Open Tabs:")
		fmt.Println()

		for _, tab := range tabs {
			if m, ok := tab.(map[string]any); ok {
				id, _ := m["id"].(string)
				tabURL, _ := m["url"].(string)
				title, _ := m["title"].(string)

				if title == "" {
					title = "(untitled)"
				}

				fmt.Printf("  [%s] %s\n", id, title)
				fmt.Printf("       %s\n", tabURL)
			}
		}
		fmt.Println()
	}
}
