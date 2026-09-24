package uaparser

import (
	"net/http"
	"strings"
	"sync"
)

// Option is a functional option for configuring a Parser instance.
type Option func(*Parser)

// WithCache sets a custom cache implementation on the Parser.
func WithCache(cache Cacher) Option {
	return func(p *Parser) {
		p.cache = cache
	}
}

// WithoutCache disables caching on the Parser.
func WithoutCache() Option {
	return func(p *Parser) {
		p.cache = NoOpCache{}
	}
}

// Parser parses User-Agent strings and Client Hints into structured results.
type Parser struct {
	cache Cacher
}

// DefaultParser is the global thread-safe Parser instance configured with a 10,000 entry sharded LRU cache.
var (
	defaultParserInstance *Parser
	defaultParserOnce     sync.Once
)

// DefaultParser returns the global singleton Parser.
func Default() *Parser {
	defaultParserOnce.Do(func() {
		defaultParserInstance = New(WithCache(NewShardedLRU(10000)))
	})
	return defaultParserInstance
}

// New creates a new Parser with the supplied options.
func New(opts ...Option) *Parser {
	p := &Parser{
		cache: NewShardedLRU(5000),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Parse parses a raw User-Agent string using the default singleton Parser.
func Parse(ua string) *Result {
	return Default().Parse(ua)
}

// ParseWithClientHints parses a User-Agent string with supplied Client Hints using the default Parser.
func ParseWithClientHints(ua string, ch ClientHints) *Result {
	return Default().ParseWithClientHints(ua, ch)
}

// ParseRequest extracts User-Agent and Client Hints from an incoming HTTP request using the default Parser.
func ParseRequest(r *http.Request) *Result {
	return Default().ParseRequest(r)
}

// ParseRequest extracts User-Agent and Client Hints from an incoming HTTP request.
func (p *Parser) ParseRequest(r *http.Request) *Result {
	if r == nil {
		return &Result{IsServer: true}
	}
	ua := r.UserAgent()
	ch := ParseClientHints(r.Header)
	return p.ParseWithClientHints(ua, ch)
}

// Parse parses a raw User-Agent string.
func (p *Parser) Parse(ua string) *Result {
	return p.ParseWithClientHints(ua, ClientHints{})
}

// ParseWithClientHints parses a User-Agent string and merges Client Hints (Sec-CH-UA-*).
func (p *Parser) ParseWithClientHints(ua string, ch ClientHints) *Result {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return &Result{IsServer: true}
	}

	// If no Client Hints, try the cache
	if p.cache != nil && isClientHintsEmpty(ch) {
		if cached, ok := p.cache.Get(ua); ok {
			return cached
		}
	}

	res := &Result{
		UA:          ua,
		ClientHints: ch,
		IsFrozen:    IsFrozenUA(ua),
		IsBot:       IsBot(ua),
	}

	// 1. Parse Browser
	for _, rule := range browserRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			res.Browser = rule.handler(matches)
			break
		}
	}

	// 2. Parse OS
	for _, rule := range osRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			res.OS = rule.handler(matches)
			break
		}
	}

	// 3. Parse Device
	for _, rule := range deviceRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			res.Device = rule.handler(matches)
			break
		}
	}

	// 4. Parse Engine
	for _, rule := range engineRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			res.Engine = rule.handler(matches)
			break
		}
	}

	// 5. Parse CPU
	for _, rule := range cpuRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			res.CPU = rule.handler(matches)
			break
		}
	}

	// 6. Handle App-Style UA formatting (e.g. Model=Redmi Note 8 Pro; Manufacturer=Xiaomi)
	appInfo := ExtractAppStyleInfo(ua)
	if appInfo.Model != "" || appInfo.Manufacturer != "" {
		if res.Device.Model == "" {
			res.Device.Model = appInfo.Model
		}
		if res.Device.Vendor == "" {
			if appInfo.Manufacturer != "" {
				res.Device.Vendor = appInfo.Manufacturer
			} else {
				res.Device.Vendor = DetectBrand(ua, appInfo.Model)
			}
		}
		if res.Device.Type == "" {
			res.Device.Type = DeviceMobile
		}
	}

	// 7. Apply Brand Detection fallback if vendor missing
	if res.Device.Vendor == "" && (res.Device.Model != "" || res.OS.Name != "") {
		if brand := DetectBrand(ua, res.Device.Model); brand != "" {
			res.Device.Vendor = brand
		}
	}

	// 8. Device type heuristics
	if res.Device.Type == "" {
		lower := strings.ToLower(ua)
		if strings.Contains(lower, "mobile") {
			res.Device.Type = DeviceMobile
		} else if strings.Contains(lower, "tablet") || strings.Contains(lower, "ipad") {
			res.Device.Type = DeviceTablet
		} else if res.OS.Name == "Windows" || res.OS.Name == "macOS" || res.OS.Name == "Linux" {
			res.Device.Type = DeviceDesktop
		}
	}

	// 9. Apply Client Hints overrides
	p.applyClientHints(res, ch)

	// 10. Check if Server UA
	res.IsServer = IsServerUA(ua, res)

	// Store in cache if plain UA
	if p.cache != nil && isClientHintsEmpty(ch) {
		p.cache.Set(ua, res)
	}

	return res
}

func (p *Parser) applyClientHints(res *Result, ch ClientHints) {
	// Mobile hint
	if ch.Mobile != nil {
		if *ch.Mobile {
			res.Device.Type = DeviceMobile
		}
	}

	// Form factors hint
	for _, factor := range ch.FormFactors {
		factorLower := strings.ToLower(factor)
		switch factorLower {
		case "tablet":
			res.Device.Type = DeviceTablet
		case "mobile":
			res.Device.Type = DeviceMobile
		case "desktop":
			res.Device.Type = DeviceDesktop
		case "automotive":
			res.Device.Type = DeviceEmbedded
		case "xr", "vr":
			res.Device.Type = DeviceXR
		}
	}

	// Model hint
	if ch.Model != "" {
		res.Device.Model = ch.Model
		if res.Device.Vendor == "" {
			res.Device.Vendor = DetectBrand(res.UA, ch.Model)
		}
	}

	// Platform & PlatformVersion hint
	if ch.Platform != "" {
		res.OS.Name = ch.Platform
		if ch.Platform == "Windows" && ch.PlatformVersion != "" {
			// Windows 11 platformVersion is >= 13.0.0
			parts := strings.Split(ch.PlatformVersion, ".")
			if len(parts) > 0 {
				if parts[0] >= "13" {
					res.OS.Version = "11"
				} else {
					res.OS.Version = "10"
				}
			}
		} else if ch.PlatformVersion != "" {
			res.OS.Version = ch.PlatformVersion
		}
	}

	// Architecture & bitness
	if ch.Architecture != "" {
		arch := strings.ToLower(ch.Architecture)
		if ch.Bitness == "64" && arch == "arm" {
			res.CPU.Architecture = "arm64"
		} else if ch.Bitness == "64" && (arch == "x86" || arch == "ia32") {
			res.CPU.Architecture = "amd64"
		} else {
			res.CPU.Architecture = arch
		}
	}

	// Brands hint (extract highest-fidelity browser name & version)
	brands := ch.FullVersionList
	if len(brands) == 0 {
		brands = ch.Brands
	}
	for _, b := range brands {
		bLower := strings.ToLower(b.Brand)
		if strings.Contains(bLower, "not") || strings.Contains(bLower, "brand") {
			continue
		}
		if strings.Contains(bLower, "edge") {
			res.Browser.Name = "Edge"
			if b.Version != "" {
				res.Browser.Version = b.Version
				res.Browser.Major = Majorize(b.Version)
			}
			break
		}
		if strings.Contains(bLower, "chrome") {
			res.Browser.Name = "Chrome"
			if b.Version != "" {
				res.Browser.Version = b.Version
				res.Browser.Major = Majorize(b.Version)
			}
		}
	}
}

func isClientHintsEmpty(ch ClientHints) bool {
	return len(ch.Brands) == 0 &&
		len(ch.FullVersionList) == 0 &&
		ch.Mobile == nil &&
		ch.Model == "" &&
		ch.Platform == "" &&
		ch.PlatformVersion == "" &&
		ch.Architecture == "" &&
		ch.Bitness == "" &&
		len(ch.FormFactors) == 0
}
