package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pinchtab/pinchtab/internal/config"
)

func printHelp() {
	fmt.Printf(`pinchtab %s - Browser control for AI agents

MODES:
  pinchtab                 Start server (default port 9867)
  pinchtab connect <name>  Get URL for a running profile instance

BROWSER COMMANDS:
  pinchtab nav <url>              Navigate to URL
  pinchtab snap [url]             Accessibility snapshot (-i, -c, -d; see SNAPSHOT FLAGS)
  pinchtab find <query> [--url u] Semantic element search (--top N, --threshold F)
  pinchtab text [url]             Extract readable text
  pinchtab screenshot [url]       Capture screenshot (--out file)
  pinchtab pdf [url]              Export PDF (see PDF FLAGS)
  pinchtab click <ref>            Click element by ref
  pinchtab type <ref> <text>      Type text into element
  pinchtab fill <ref> <text>      Fill/replace input value
  pinchtab press <key>            Press keyboard key (Enter, Tab, Escape, ...)
  pinchtab hover <ref>            Hover over element
  pinchtab scroll [ref]           Scroll element into view (or page)
  pinchtab select <ref> <value>   Select dropdown option
  pinchtab eval <expression>      Evaluate JavaScript

INSTANCE COMMANDS:
  pinchtab launch [profile]  Launch instance (--headed, --port N)
  pinchtab stop <id>         Stop a running instance
  pinchtab instances         List running instances

TAB COMMANDS:
  pinchtab open [url]        Open new tab (optional URL)
  pinchtab close <id>        Close tab by ID
  pinchtab tabs              List open tabs

MANAGEMENT COMMANDS:
  pinchtab health        Server health check
  pinchtab profiles      List available profiles
  pinchtab config init   Initialize config file
  pinchtab config show   Display current configuration
  pinchtab help          Show this help

SNAPSHOT FLAGS:
  -i, --interactive    Interactive elements only
  -c, --compact        Compact plain-text snapshot output
  -d, --diff           Only changes since last snapshot
  -s, --selector CSS   Scope to CSS selector
  --max-tokens N       Truncate to ~N tokens
  --depth N            Max tree depth
  --tab ID             Target specific tab

PDF FLAGS:
  -o, --out FILE             Save PDF locally (default: page.pdf)
  --landscape                Landscape orientation
  --tab ID                   Target specific tab
  --paper-width N            Paper width (inches)
  --paper-height N           Paper height (inches)
  --margin-top N             Top margin (inches)
  --margin-bottom N          Bottom margin (inches)
  --margin-left N            Left margin (inches)
  --margin-right N           Right margin (inches)
  --scale N                  Print scale 0.1-2.0
  --page-ranges RANGE        Pages to export (e.g. 1-3,5)
  --prefer-css-page-size     Honor CSS @page size
  --display-header-footer    Show header/footer
  --header-template HTML     Header template
  --footer-template HTML     Footer template
  --generate-tagged-pdf      Generate accessible/tagged PDF
  --generate-document-outline  Embed document outline
  --file-output              Save PDF in server state dir
  --path PATH                Custom server-side path (with --file-output)

ENVIRONMENT:
  PINCHTAB_URL    Server URL (default: http://127.0.0.1:9867)
  PINCHTAB_TOKEN  Auth token for API requests
`, version)
}

var cliCommands = map[string]bool{
	// Management
	"health":   true,
	"help":     true,
	"config":   true,
	"profiles": true,
	"connect":  true,

	// Instances
	"launch":    true,
	"stop":      true,
	"instances": true,

	// Tabs
	"open":  true,
	"close": true,
	"tabs":  true,

	// Browser: navigation & read
	"nav":        true,
	"navigate":   true,
	"snap":       true,
	"snapshot":   true,
	"find":       true,
	"text":       true,
	"screenshot": true,
	"ss":         true,
	"pdf":        true,

	// Browser: actions
	"click":  true,
	"type":   true,
	"fill":   true,
	"press":  true,
	"hover":  true,
	"scroll": true,
	"select": true,

	// Browser: eval
	"eval":     true,
	"evaluate": true,
}

func isCLICommand(cmd string) bool {
	return cliCommands[cmd]
}

func runCLI(cfg *config.RuntimeConfig) {
	cmd := os.Args[1]

	base := fmt.Sprintf("http://%s:%s", cfg.Bind, cfg.Port)
	if envURL := os.Getenv("PINCHTAB_URL"); envURL != "" {
		base = strings.TrimRight(envURL, "/")
	}

	token := cfg.Token
	if envToken := os.Getenv("PINCHTAB_TOKEN"); envToken != "" {
		token = envToken
	}

	client := &http.Client{Timeout: 30 * time.Second}

	args := os.Args[2:]

	switch cmd {
	// Management
	case "health":
		cliHealth(client, base, token)
	case "profiles":
		cliProfiles(client, base, token)
	case "help":
		printHelp()
	case "config":
		// Handled in main.go

	// Instances
	case "launch":
		cliLaunch(client, base, token, args)
	case "stop":
		cliStop(client, base, token, args)
	case "instances":
		cliInstances(client, base, token)

	// Tabs
	case "open":
		cliOpen(client, base, token, args)
	case "close":
		cliClose(client, base, token, args)
	case "tabs":
		cliTabs(client, base, token)

	// Browser: navigation & read
	case "nav", "navigate":
		cliNavigate(client, base, token, args)
	case "snap", "snapshot":
		cliSnapshot(client, base, token, args)
	case "find":
		cliFind(client, base, token, args)
	case "text":
		cliText(client, base, token, args)
	case "screenshot", "ss":
		cliScreenshot(client, base, token, args)
	case "pdf":
		cliPDF(client, base, token, args)
	case "eval", "evaluate":
		cliEval(client, base, token, args)

	// Browser: actions
	case "click", "hover", "scroll", "press", "type", "fill", "select":
		cliAction(client, base, token, cmd, args)
	}
}
