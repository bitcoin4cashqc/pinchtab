package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func cliNavigate(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: pinchtab nav <url>")
		os.Exit(1)
	}
	result := doPost(client, base, token, "/navigate", map[string]any{"url": args[0]})
	if tabID, ok := result["tabId"].(string); ok {
		fmt.Printf("Navigated [%s] → %s\n", tabID, args[0])
	}
}

func cliSnapshot(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	rawFormat := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-i", "--interactive":
			params.Set("filter", "interactive")
		case "-c", "--compact":
			params.Set("format", "compact")
			rawFormat = "compact"
		case "--text":
			params.Set("format", "text")
			rawFormat = "text"
		case "-d", "--diff":
			params.Set("diff", "true")
		case "--depth":
			if i+1 < len(args) {
				i++
				params.Set("depth", args[i])
			}
		case "-s", "--selector":
			if i+1 < len(args) {
				i++
				params.Set("selector", args[i])
			}
		case "--max-tokens":
			if i+1 < len(args) {
				i++
				params.Set("maxTokens", args[i])
			}
		case "--tab":
			if i+1 < len(args) {
				i++
				params.Set("tabId", args[i])
			}
		default:
			if strings.HasPrefix(args[i], "http") {
				params.Set("url", args[i])
			}
		}
	}

	if rawFormat != "" && params.Get("diff") != "true" {
		rawBody := doGetRaw(client, base, token, "/snapshot", params)
		fmt.Println(string(rawBody))
		return
	}

	result := doGet(client, base, token, "/snapshot", params)
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}

func cliFind(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: pinchtab find <query> [--url <url>] [--top N] [--threshold F]")
		os.Exit(1)
	}

	body := map[string]any{"query": args[0]}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--url":
			if i+1 < len(args) {
				i++
				body["url"] = args[i]
			}
		case "--top":
			if i+1 < len(args) {
				i++
				if n, err := strconv.Atoi(args[i]); err == nil {
					body["topK"] = n
				}
			}
		case "--threshold":
			if i+1 < len(args) {
				i++
				if f, err := strconv.ParseFloat(args[i], 64); err == nil {
					body["threshold"] = f
				}
			}
		}
	}

	result := doPost(client, base, token, "/find", body)

	if matches, ok := result["matches"].([]any); ok {
		if len(matches) == 0 {
			fmt.Println("No matches found")
			return
		}
		for _, m := range matches {
			if entry, ok := m.(map[string]any); ok {
				ref, _ := entry["ref"].(string)
				name, _ := entry["name"].(string)
				role, _ := entry["role"].(string)
				score, _ := entry["score"].(float64)
				fmt.Printf("  [%s] %.2f  %s: %s\n", ref, score, role, name)
			}
		}
		if best, ok := result["best_ref"].(string); ok {
			fmt.Printf("\nBest: %s\n", best)
		}
	}
}

func cliText(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "http") {
			params.Set("url", arg)
		}
	}
	result := doGet(client, base, token, "/text", params)
	if text, ok := result["text"].(string); ok {
		fmt.Println(text)
	} else {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	}
}

func cliScreenshot(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	outFile := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--out":
			if i+1 < len(args) {
				i++
				outFile = args[i]
			}
		default:
			if strings.HasPrefix(args[i], "http") {
				params.Set("url", args[i])
			}
		}
	}
	if outFile == "" {
		outFile = "screenshot.jpg"
	}
	params.Set("output", "raw")

	rawBody := doGetRaw(client, base, token, "/screenshot", params)
	if err := os.WriteFile(outFile, rawBody, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Screenshot saved to %s (%d bytes)\n", outFile, len(rawBody))
}

func cliPDF(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	outFile := ""
	serverSideOutput := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o", "--out":
			if i+1 < len(args) {
				i++
				outFile = args[i]
			}
		case "--landscape":
			params.Set("landscape", "true")
		case "--tab":
			if i+1 < len(args) {
				i++
				params.Set("tabId", args[i])
			}
		case "--paper-width":
			if i+1 < len(args) {
				i++
				params.Set("paperWidth", args[i])
			}
		case "--paper-height":
			if i+1 < len(args) {
				i++
				params.Set("paperHeight", args[i])
			}
		case "--margin-top":
			if i+1 < len(args) {
				i++
				params.Set("marginTop", args[i])
			}
		case "--margin-bottom":
			if i+1 < len(args) {
				i++
				params.Set("marginBottom", args[i])
			}
		case "--margin-left":
			if i+1 < len(args) {
				i++
				params.Set("marginLeft", args[i])
			}
		case "--margin-right":
			if i+1 < len(args) {
				i++
				params.Set("marginRight", args[i])
			}
		case "--scale":
			if i+1 < len(args) {
				i++
				params.Set("scale", args[i])
			}
		case "--page-ranges":
			if i+1 < len(args) {
				i++
				params.Set("pageRanges", args[i])
			}
		case "--prefer-css-page-size":
			params.Set("preferCSSPageSize", "true")
		case "--display-header-footer":
			params.Set("displayHeaderFooter", "true")
		case "--header-template":
			if i+1 < len(args) {
				i++
				params.Set("headerTemplate", args[i])
			}
		case "--footer-template":
			if i+1 < len(args) {
				i++
				params.Set("footerTemplate", args[i])
			}
		case "--generate-tagged-pdf":
			params.Set("generateTaggedPDF", "true")
		case "--generate-document-outline":
			params.Set("generateDocumentOutline", "true")
		case "--file-output":
			serverSideOutput = true
		case "--path":
			if i+1 < len(args) {
				i++
				params.Set("path", args[i])
			}
		default:
			if strings.HasPrefix(args[i], "http") {
				params.Set("url", args[i])
			}
		}
	}

	if serverSideOutput {
		if outFile != "" {
			fmt.Fprintln(os.Stderr, "Cannot combine --file-output with --out")
			os.Exit(1)
		}
		params.Set("output", "file")
		result := doGet(client, base, token, "/pdf", params)
		if path, ok := result["path"].(string); ok {
			fmt.Printf("PDF saved on server: %s\n", path)
			return
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return
	}

	if outFile == "" {
		outFile = "page.pdf"
	}
	params.Set("raw", "true")

	rawBody := doGetRaw(client, base, token, "/pdf", params)
	if err := os.WriteFile(outFile, rawBody, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("PDF saved to %s (%d bytes)\n", outFile, len(rawBody))
}

func cliAction(client *http.Client, base, token, kind string, args []string) {
	body := map[string]any{"kind": kind}

	switch kind {
	case "click", "hover", "focus":
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Usage: pinchtab %s <ref>\n", kind)
			os.Exit(1)
		}
		body["ref"] = args[0]
	case "type":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pinchtab type <ref> <text>")
			os.Exit(1)
		}
		body["ref"] = args[0]
		body["text"] = strings.Join(args[1:], " ")
	case "fill":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pinchtab fill <ref> <text>")
			os.Exit(1)
		}
		body["ref"] = args[0]
		body["text"] = strings.Join(args[1:], " ")
	case "press":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: pinchtab press <key>")
			os.Exit(1)
		}
		body["key"] = args[0]
	case "scroll":
		if len(args) > 0 {
			body["ref"] = args[0]
		}
	case "select":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pinchtab select <ref> <value>")
			os.Exit(1)
		}
		body["ref"] = args[0]
		body["value"] = args[1]
	}

	result := doPost(client, base, token, "/action", body)
	if msg, ok := result["status"].(string); ok {
		fmt.Println(msg)
	} else {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	}
}

func cliEval(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: pinchtab eval <expression>")
		os.Exit(1)
	}
	expr := strings.Join(args, " ")
	result := doPost(client, base, token, "/evaluate", map[string]any{"expression": expr})
	if val, ok := result["result"]; ok {
		out, _ := json.MarshalIndent(val, "", "  ")
		fmt.Println(string(out))
	}
}
