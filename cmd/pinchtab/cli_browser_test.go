package main

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestCLISnapshotFlags(t *testing.T) {
	m := newMockServer()
	m.response = `{"ok":true}`
	defer m.close()
	client := m.server.Client()

	cliSnapshot(client, m.base(), "", []string{
		"-i", "-c", "-d",
		"--depth", "3",
		"--max-tokens", "100",
		"--selector", ".btn",
		"--tab", "tab_123",
		"https://example.com",
	})

	if m.lastPath != "/snapshot" {
		t.Fatalf("expected /snapshot, got %s", m.lastPath)
	}
	q, err := url.ParseQuery(m.lastQuery)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("filter") != "interactive" {
		t.Fatalf("expected filter=interactive, got %q", q.Get("filter"))
	}
	if q.Get("format") != "compact" {
		t.Fatalf("expected format=compact, got %q", q.Get("format"))
	}
	if q.Get("diff") != "true" {
		t.Fatalf("expected diff=true, got %q", q.Get("diff"))
	}
	if q.Get("depth") != "3" {
		t.Fatalf("expected depth=3, got %q", q.Get("depth"))
	}
	if q.Get("maxTokens") != "100" {
		t.Fatalf("expected maxTokens=100, got %q", q.Get("maxTokens"))
	}
	if q.Get("selector") != ".btn" {
		t.Fatalf("expected selector=.btn, got %q", q.Get("selector"))
	}
	if q.Get("tabId") != "tab_123" {
		t.Fatalf("expected tabId=tab_123, got %q", q.Get("tabId"))
	}
	if q.Get("url") != "https://example.com" {
		t.Fatalf("expected url=https://example.com, got %q", q.Get("url"))
	}
}

func TestCLIPDFFlags_LocalOutput(t *testing.T) {
	m := newMockServer()
	m.response = "%PDF-fake"
	defer m.close()
	client := m.server.Client()

	tmp := t.TempDir()
	outFile := filepath.Join(tmp, "test.pdf")

	cliPDF(client, m.base(), "", []string{
		"https://example.com",
		"-o", outFile,
		"--tab", "tab_abc",
		"--landscape",
		"--paper-width", "8.5",
		"--paper-height", "11",
		"--margin-top", "0.1",
		"--margin-bottom", "0.2",
		"--margin-left", "0.3",
		"--margin-right", "0.4",
		"--scale", "1.1",
		"--page-ranges", "1-2",
		"--prefer-css-page-size",
		"--display-header-footer",
		"--header-template", "<h1>x</h1>",
		"--footer-template", "<p>y</p>",
		"--generate-tagged-pdf",
		"--generate-document-outline",
	})

	if m.lastPath != "/pdf" {
		t.Fatalf("expected /pdf, got %s", m.lastPath)
	}
	q, err := url.ParseQuery(m.lastQuery)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("raw") != "true" {
		t.Fatalf("expected raw=true, got %q", q.Get("raw"))
	}
	if q.Get("tabId") != "tab_abc" {
		t.Fatalf("expected tabId=tab_abc, got %q", q.Get("tabId"))
	}
	if q.Get("landscape") != "true" {
		t.Fatalf("expected landscape=true, got %q", q.Get("landscape"))
	}
	if q.Get("paperWidth") != "8.5" || q.Get("paperHeight") != "11" {
		t.Fatalf("expected paper size params, got %q x %q", q.Get("paperWidth"), q.Get("paperHeight"))
	}
	if q.Get("marginTop") != "0.1" || q.Get("marginBottom") != "0.2" || q.Get("marginLeft") != "0.3" || q.Get("marginRight") != "0.4" {
		t.Fatalf("unexpected margin params: %s", m.lastQuery)
	}
	if q.Get("scale") != "1.1" || q.Get("pageRanges") != "1-2" {
		t.Fatalf("unexpected scale/pageRanges params: %s", m.lastQuery)
	}
	if q.Get("preferCSSPageSize") != "true" || q.Get("displayHeaderFooter") != "true" {
		t.Fatalf("expected css/header flags, got %s", m.lastQuery)
	}
	if q.Get("headerTemplate") != "<h1>x</h1>" || q.Get("footerTemplate") != "<p>y</p>" {
		t.Fatalf("expected templates in query, got %s", m.lastQuery)
	}
	if q.Get("generateTaggedPDF") != "true" || q.Get("generateDocumentOutline") != "true" {
		t.Fatalf("expected tagged/outline flags, got %s", m.lastQuery)
	}

	if _, err := os.Stat(outFile); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
}

func TestCLIPDFFlags_ServerOutput(t *testing.T) {
	m := newMockServer()
	m.response = `{"path":"/tmp/server.pdf","size":1234}`
	defer m.close()
	client := m.server.Client()

	cliPDF(client, m.base(), "", []string{
		"--file-output",
		"--path", "pdfs/output.pdf",
	})

	if m.lastPath != "/pdf" {
		t.Fatalf("expected /pdf, got %s", m.lastPath)
	}
	q, err := url.ParseQuery(m.lastQuery)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("output") != "file" {
		t.Fatalf("expected output=file, got %q", q.Get("output"))
	}
	if q.Get("path") != "pdfs/output.pdf" {
		t.Fatalf("expected path=pdfs/output.pdf, got %q", q.Get("path"))
	}
	if q.Get("raw") != "" {
		t.Fatalf("expected no raw param for server output, got %q", q.Get("raw"))
	}
}
