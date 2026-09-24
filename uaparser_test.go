package uaparser

import (
	"net/http"
	"testing"
)

func TestParseBrowsers(t *testing.T) {
	tests := []struct {
		name         string
		ua           string
		expectName   string
		expectType   BrowserType
	}{
		{
			name:       "Google Chrome on Mac",
			ua:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
			expectName: "Chrome",
		},
		{
			name:       "Mobile Chrome on Android",
			ua:         "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
			expectName: "Chrome",
		},
		{
			name:       "Safari on iPhone",
			ua:         "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1",
			expectName: "Mobile Safari",
		},
		{
			name:       "Firefox on Windows",
			ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0",
			expectName: "Firefox",
		},
		{
			name:       "Microsoft Edge",
			ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
			expectName: "Edge",
		},
		{
			name:       "Edge WebView InApp",
			ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64; WebView/3.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edge/124.0.0.0",
			expectName: "Edge WebView",
			expectType: BrowserTypeInApp,
		},
		{
			name:       "Samsung Internet",
			ua:         "Mozilla/5.0 (Linux; Android 14; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/24.0 Chrome/117.0.0.0 Mobile Safari/537.36",
			expectName: "Samsung Internet",
		},
		{
			name:       "WeChat InApp",
			ua:         "Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/111.0.5563.116 Mobile Safari/537.36 MicroMessenger/8.0.38.2400(0x28002635)",
			expectName: "WeChat",
			expectType: BrowserTypeInApp,
		},
		{
			name:       "Chrome Headless Crawler",
			ua:         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/124.0.6367.60 Safari/537.36",
			expectName: "Chrome Headless",
			expectType: BrowserTypeCrawler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Parse(tt.ua)
			if res.Browser.Name != tt.expectName {
				t.Errorf("expected Browser Name %q, got %q", tt.expectName, res.Browser.Name)
			}
			if tt.expectType != "" && res.Browser.Type != tt.expectType {
				t.Errorf("expected Browser Type %q, got %q", tt.expectType, res.Browser.Type)
			}
		})
	}
}

func TestParseOS(t *testing.T) {
	tests := []struct {
		name          string
		ua            string
		expectOS      string
		expectVersion string
	}{
		{
			name:          "Windows 10",
			ua:            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
			expectOS:      "Windows",
			expectVersion: "10",
		},
		{
			name:          "macOS Sonoma",
			ua:            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			expectOS:      "macOS",
			expectVersion: "10.15.7",
		},
		{
			name:          "iOS on iPhone",
			ua:            "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15",
			expectOS:      "iOS",
			expectVersion: "17.4.1",
		},
		{
			name:          "Android 14",
			ua:            "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36",
			expectOS:      "Android",
			expectVersion: "14",
		},
		{
			name:     "Ubuntu Linux",
			ua:       "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
			expectOS: "Ubuntu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Parse(tt.ua)
			if res.OS.Name != tt.expectOS {
				t.Errorf("expected OS %q, got %q", tt.expectOS, res.OS.Name)
			}
			if tt.expectVersion != "" && res.OS.Version != tt.expectVersion {
				t.Errorf("expected OS Version %q, got %q", tt.expectVersion, res.OS.Version)
			}
		})
	}
}

func TestParseDevices(t *testing.T) {
	tests := []struct {
		name         string
		ua           string
		expectVendor string
		expectModel  string
		expectType   DeviceType
	}{
		{
			name:         "Apple iPad",
			ua:           "Mozilla/5.0 (iPad; CPU OS 17_4 like Mac OS X) AppleWebKit/605.1.15",
			expectVendor: "Apple",
			expectModel:  "iPad",
			expectType:   DeviceTablet,
		},
		{
			name:         "Apple iPhone",
			ua:           "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15",
			expectVendor: "Apple",
			expectModel:  "iPhone",
			expectType:   DeviceMobile,
		},
		{
			name:         "Samsung Galaxy Tab (SM-T870)",
			ua:           "Mozilla/5.0 (Linux; Android 12; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36",
			expectVendor: "Samsung",
			expectModel:  "SM-T870",
			expectType:   DeviceTablet,
		},
		{
			name:         "Samsung Galaxy S23 (SM-S918B)",
			ua:           "Mozilla/5.0 (Linux; Android 14; SM-S918B) AppleWebKit/537.36 Mobile Safari/537.36",
			expectVendor: "Samsung",
			expectModel:  "SM-S918B",
			expectType:   DeviceMobile,
		},
		{
			name:         "LG Mobile Phone",
			ua:           "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30",
			expectVendor: "LG",
			expectModel:  "LG-L160L",
			expectType:   DeviceMobile,
		},
		{
			name:         "Google Pixel 8",
			ua:           "Mozilla/5.0 (Linux; Android 14; Pixel 8 Build/UD1A.230803.041) AppleWebKit/537.36 Mobile Safari/537.36",
			expectVendor: "Google",
			expectModel:  "Pixel 8",
			expectType:   DeviceMobile,
		},
		{
			name:         "App-style Xiaomi UA",
			ua:           "CustomApp/2.0 (Model=Redmi Note 8 Pro; Manufacturer=Xiaomi) Android/11",
			expectVendor: "Xiaomi",
			expectModel:  "Redmi Note 8 Pro",
			expectType:   DeviceMobile,
		},
		{
			name:         "PlayStation 5 Console",
			ua:           "Mozilla/5.0 (PlayStation 5 7.00) AppleWebKit/605.1.15 (KHTML, like Gecko)",
			expectVendor: "Sony",
			expectType:   DeviceConsole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Parse(tt.ua)
			if tt.expectVendor != "" && res.Device.Vendor != tt.expectVendor {
				t.Errorf("expected Vendor %q, got %q", tt.expectVendor, res.Device.Vendor)
			}
			if tt.expectModel != "" && res.Device.Model != tt.expectModel {
				t.Errorf("expected Model %q, got %q", tt.expectModel, res.Device.Model)
			}
			if tt.expectType != "" && res.Device.Type != tt.expectType {
				t.Errorf("expected Device Type %q, got %q", tt.expectType, res.Device.Type)
			}
		})
	}
}

func TestClientHints(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Full-Version-List", `"Chromium";v="124.0.6367.60", "Google Chrome";v="124.0.6367.60"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-CH-UA-Platform-Version", `"15.0.0"`) // Windows 11
	req.Header.Set("Sec-CH-UA-Arch", `"x86"`)
	req.Header.Set("Sec-CH-UA-Bitness", `"64"`)

	res := ParseRequest(req)
	if res.Browser.Name != "Chrome" {
		t.Errorf("expected Chrome, got %s", res.Browser.Name)
	}
	if res.Browser.Version != "124.0.6367.60" {
		t.Errorf("expected 124.0.6367.60, got %s", res.Browser.Version)
	}
	if res.OS.Name != "Windows" {
		t.Errorf("expected Windows, got %s", res.OS.Name)
	}
	if res.OS.Version != "11" {
		t.Errorf("expected Windows 11 via platformVersion, got %s", res.OS.Version)
	}
	if res.CPU.Architecture != "amd64" {
		t.Errorf("expected amd64 architecture, got %s", res.CPU.Architecture)
	}
}

func TestBotAndCrawlerDetection(t *testing.T) {
	tests := []struct {
		ua            string
		isBot         bool
		isCrawler     bool
		isAICrawler   bool
		isAIAssistant bool
	}{
		{
			ua:          "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			isBot:       true,
			isCrawler:   true,
			isAICrawler: false,
		},
		{
			ua:        "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b",
			isBot:     true,
			isCrawler: true,
		},
		{
			ua:        "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
			isBot:     true,
			isCrawler: true,
		},
		{
			ua:        "Screaming Frog SEO Spider/19.0",
			isBot:     true,
			isCrawler: true,
		},
		{
			ua:          "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.2; +https://openai.com/gptbot)",
			isBot:       true,
			isCrawler:   true,
			isAICrawler: true,
		},
		{
			ua:          "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; ClaudeBot/1.0; +claudebot@anthropic.com)",
			isBot:       true,
			isCrawler:   true,
			isAICrawler: true,
		},
		{
			ua:            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 (ChatGPT-User)",
			isBot:         true,
			isCrawler:     false,
			isAIAssistant: true,
		},
		{
			ua:        "curl/7.88.1",
			isBot:     false,
			isCrawler: false,
		},
		{
			ua:        "Go-http-client/1.1",
			isBot:     false,
			isCrawler: false,
		},
		{
			ua:        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
			isBot:     false,
			isCrawler: false,
		},
	}

	for _, tt := range tests {
		if got := IsBot(tt.ua); got != tt.isBot {
			t.Errorf("UA %q: expected IsBot=%v, got %v", tt.ua, tt.isBot, got)
		}
		if got := IsCrawler(tt.ua); got != tt.isCrawler {
			t.Errorf("UA %q: expected IsCrawler=%v, got %v", tt.ua, tt.isCrawler, got)
		}
		if tt.isAICrawler && !IsAICrawler(tt.ua) {
			t.Errorf("UA %q: expected IsAICrawler=true", tt.ua)
		}
		if tt.isAIAssistant && !IsAIAssistant(tt.ua) {
			t.Errorf("UA %q: expected IsAIAssistant=true", tt.ua)
		}
	}
}

func TestFrozenUA(t *testing.T) {
	frozen := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	if !IsFrozenUA(frozen) {
		t.Fatalf("expected frozen UA to be true for %s", frozen)
	}

	nonFrozen := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.60 Safari/537.36"
	if IsFrozenUA(nonFrozen) {
		t.Fatalf("expected non-frozen UA for non-reduced Chrome string")
	}
}

func TestEuropeTimezones(t *testing.T) {
	if !IsFromEU("Europe/Berlin") {
		t.Errorf("Europe/Berlin should be EU")
	}
	if !IsFromEU("Europe/Paris") {
		t.Errorf("Europe/Paris should be EU")
	}
	if IsFromEU("America/New_York") {
		t.Errorf("America/New_York should not be EU")
	}

	if !IsFromEEA("Europe/Oslo") {
		t.Errorf("Europe/Oslo should be EEA")
	}
	if !IsFromEFTA("Europe/Zurich") {
		t.Errorf("Europe/Zurich should be EFTA")
	}
	if !IsFromSchengen("Europe/Stockholm") {
		t.Errorf("Europe/Stockholm should be Schengen")
	}
}

func BenchmarkParseParallelWithoutCache(b *testing.B) {
	p := New(WithoutCache(), WithParallel(true))
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	b.ResetTimer()
	for b.Loop() {
		_ = p.Parse(ua)
	}
}

func BenchmarkParseSequentialWithoutCache(b *testing.B) {
	p := New(WithoutCache(), WithParallel(false))
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	b.ResetTimer()
	for b.Loop() {
		_ = p.Parse(ua)
	}
}

func BenchmarkParseWithCache(b *testing.B) {
	p := New()
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	// Warm cache
	_ = p.Parse(ua)

	b.ResetTimer()
	for b.Loop() {
		_ = p.Parse(ua)
	}
}
