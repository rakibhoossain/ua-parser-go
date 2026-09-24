# Maintenance & Rule Update Guide

`ua-parser-go` is designed with high standards of maintainability, separation of concerns, and clean engineering. As new browsers, mobile hardware, AI bots, and referral sources emerge, follow this guide to keep the library updated and reliable.

---

## 1. Updating the Referrer Database

The referral detection engine uses an embedded dictionary of 2,500+ domains mapped to semantic categories (`search`, `social`, `ai`, `tech`, `commerce`, `email`, `content`).

### 1.1 Automatic Upstream Sync
To refresh the referral database from the upstream Snowplow repository and merge curated AI/tech entries:

```bash
# Using go generate
go generate ./...

# Or directly running the generator CLI
go run ./cmd/generator
```

This script:
1. Fetches the latest `referers-latest.json` from the Snowplow asset repository.
2. Merges modern AI agents, developer platforms, and social networks defined in `cmd/generator/main.go:extraReferrers`.
3. Validates and formats `referrer/data_generated.go` deterministically.

### 1.2 Adding Custom Referrers
If a new platform or search engine needs to be added:
1. Open `cmd/generator/main.go`.
2. Add the domain and metadata to `extraReferrers`:
   ```go
   "newplatform.ai": {Type: "ai", Name: "NewPlatform"},
   ```
3. Run `go generate ./...`.
4. Add a test case in `referrer/referrer_test.go` and run `go test ./...`.

---

## 2. Updating User-Agent & Device Detection Rules

### 2.1 Adding a New Browser
Open `regexes.go` and append to `initBrowserRules()`:
```go
{
    regex: regexp.MustCompile(`(?i)newbrowser/([\w\.]+)`),
    handler: func(m []string) Browser {
        return Browser{Name: "NewBrowser", Version: m[1], Major: Majorize(m[1])}
    },
},
```
*Tip: Place more specific browser patterns before general ones (e.g. Edge and Opera before Chromium/Chrome).*

### 2.2 Adding New Device or Mobile Hardware
Open `regexes.go` and append to `initDeviceRules()`:
```go
{
    regex: regexp.MustCompile(`(?i)\b(myphone\s*[0-9]+)\b`),
    handler: func(m []string) Device {
        return Device{Vendor: "BrandName", Model: m[1], Type: DeviceMobile}
    },
},
```
OpenPanel app-style patterns (`Model=...; Manufacturer=...`) and brand normalization dictionaries in `helpers.go:brandPatterns` are automatically applied as fallbacks.

### 2.3 Adding New AI Crawlers or Assistants
Open `bots.go` and add the bot token to the corresponding regex:
- **`aiAssistantsRegex`**: For interactive user-facing browsing (e.g. `ChatGPT-User`, `Claude-Web`).
- **`aiCrawlersRegex`**: For data scrapers and LLM training spiders (e.g. `GPTBot`, `ClaudeBot`, `PerplexityBot`).

---

## 3. Testing and Verification Workflow

Before submitting a pull request or tagging a release:

```bash
# 1. Run full test suite
go test -v ./...

# 2. Run benchmarks to verify zero regressions
go test -v -bench=. ./...

# 3. Check formatting and vet
go fmt ./...
go vet ./...
```
