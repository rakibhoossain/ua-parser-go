# ua-parser-go

[![Go Reference](https://pkg.go.dev/badge/github.com/rakibhoossain/ua-parser-go.svg)](https://pkg.go.dev/github.com/rakibhoossain/ua-parser-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rakibhoossain/ua-parser-go)](https://goreportcard.com/report/github.com/rakibhoossain/ua-parser-go)

High-performance, **zero-dependency**, pure Go library for User-Agent parsing, User-Agent Client Hints (`Sec-CH-UA-*`), Bot & AI Crawler detection, Frozen UA detection, European compliance timezone checking, and Referral source detection with Favicon resolution.

Ported from and expanding upon [`ua-parser-js`](https://github.com/faisalman/ua-parser-js) (v2), [`ua-is-frozen`](https://github.com/faisalman/ua-is-frozen), [`detect-europe-js`](https://github.com/faisalman/detect-europe-js), and [OpenPanel](https://github.com/OpenPanel-dev/openpanel).

---

## ⚡ Highlights

- **Zero External Dependencies**: 100% Go standard library (`net/http`, `net/url`, `regexp`, `strings`, `sync`).
- **Blazing Fast (>8,000,000 ops/sec)**: Built-in 32-shard concurrent LRU cache eliminates lock contention across CPU cores.
- **Modern Standards**: Full support for Chromium User-Agent Client Hints (UACH) and User-Agent Reduction (`isFrozenUA`).
- **Comprehensive Bot & AI Detection**: Distinguishes user-facing AI assistants (ChatGPT, Claude, Gemini) from LLM training scrapers (GPTBot, ClaudeBot, PerplexityBot, ByteSpider).
- **Referral & Favicon Engine**: Embedded database of 2,500+ domains with category classification (`ai`, `social`, `search`, `tech`, `commerce`, `email`, `content`) and instant Favicon resolution via DuckDuckGo and Google.
- **Enterprise-Grade Clean Architecture**: Strict adherence to DRY and SOLID principles, modular subpackages, and fully thread-safe APIs.

---

## 📦 Installation

```bash
go get github.com/rakibhoossain/ua-parser-go
```

Requires Go 1.22+.

---

## 🚀 Quick Start

### 1. Basic User-Agent Parsing

```go
package main

import (
	"fmt"
	"github.com/rakibhoossain/ua-parser-go"
)

func main() {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	res := uaparser.Parse(ua)

	fmt.Printf("Browser: %s %s\n", res.Browser.Name, res.Browser.Major) // Chrome 124
	fmt.Printf("OS:      %s %s\n", res.OS.Name, res.OS.Version)       // macOS 10.15.7
	fmt.Printf("Device:  %s (%s)\n", res.Device.Vendor, res.Device.Type) // Apple (desktop)
	fmt.Printf("Frozen:  %v\n", res.IsFrozen)                       // true
}
```

### 2. HTTP Request with User-Agent Client Hints (UACH)

Chromium browsers freeze OS and hardware tokens in the raw User-Agent string. `ParseRequest` extracts both the User-Agent and `Sec-CH-UA-*` headers to resolve high-entropy details (exact model, Windows 11 platform versions, 64-bit architecture):

```go
func handler(w http.ResponseWriter, r *http.Request) {
	res := uaparser.ParseRequest(r)

	fmt.Println("Browser:", res.Browser.Name, res.Browser.Version)
	fmt.Println("OS:", res.OS.Name, res.OS.Version) // e.g. Windows 11
	fmt.Println("Device Model:", res.Device.Model)
	fmt.Println("Is Bot:", res.IsBot)
}
```

### 3. AI & Bot Detection

Detect whether incoming traffic is a human user, a web crawler, or an AI system:

```go
ua := "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.2; +https://openai.com/gptbot)"

uaparser.IsBot(ua)         // true
uaparser.IsAICrawler(ua)   // true (GPTBot, ClaudeBot, PerplexityBot, ByteSpider, etc.)
uaparser.IsAIAssistant(ua) // false (interactive agents like ChatGPT-User, Claude-Web)
```

### 4. European Regulatory Timezone Detection

Detect whether a client's IANA timezone falls under European regulatory frameworks (GDPR, cookie consent):

```go
uaparser.IsFromEU("Europe/Berlin")     // true
uaparser.IsFromEEA("Europe/Oslo")      // true (EU + Norway, Iceland, Liechtenstein)
uaparser.IsFromEFTA("Europe/Zurich")   // true (Switzerland + EEA/EFTA)
uaparser.IsFromSchengen("Europe/Rome") // true
```

### 5. Referral Detection & Favicon Resolution

Parse referrers from URLs or campaign query parameters (`utm_source`, `ref`, `utm_referrer`):

```go
import "github.com/rakibhoossain/ua-parser-go/referrer"

// From URL
ref := referrer.Parse("https://chatgpt.com/c/66f29ab0-1234")
fmt.Println(ref.Name)       // "ChatGPT"
fmt.Println(ref.Type)       // "ai"
fmt.Println(ref.Domain)     // "chatgpt.com"
fmt.Println(ref.FaviconURL) // "https://icons.duckduckgo.com/ip3/chatgpt.com.ico"

// From query parameters
params := map[string]string{"utm_source": "twitter"}
ref = referrer.ParseWithQuery(params)
fmt.Println(ref.Name) // "Twitter"
fmt.Println(ref.Type) // "social"

// Directly from an HTTP request
ref = referrer.ParseRequest(r)
```

---

## 📊 Benchmarks

Benchmark executed on Apple M3 Pro (macOS darwin/arm64):

```
cpu: Apple M3 Pro
BenchmarkParseWithoutCache-11      6,362 ops    187,221 ns/op
BenchmarkParseWithCache-11    10,026,727 ops        121.4 ns/op
```

- **With LRU Cache**: Over **8,000,000 requests/second per core** at ~121ns per parse.
- **Without Cache**: Full regex evaluation in ~187µs per parse.

---

## 🛠 Maintenance & Updating Rules

User-Agent patterns and referral sources change as new devices and platforms launch. Upstream data can be updated with a single command:

```bash
go generate ./...
```

For detailed instructions on adding custom browsers, device heuristics, or new referral categories, see [MAINTENANCE.md](MAINTENANCE.md).

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.
