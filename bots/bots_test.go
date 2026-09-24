package bots

import (
	"testing"
)

func TestDetectBot(t *testing.T) {
	t.Run("does not flag legitimate runtime and client user agents", func(t *testing.T) {
		allowlisted := []string{"node", "Node.js"}
		backendClients := []string{
			"undici",
			"axios/1.6.0",
			"node-fetch/2.6.7",
			"got (https://github.com/sindresorhus/got)",
			"PostmanRuntime/7.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		}

		for _, ua := range append(allowlisted, backendClients...) {
			match := DetectBot(ua)
			if match != nil {
				t.Errorf("expected %q not to be detected as a bot, got %+v", ua, match)
			}
			if IsBot(ua) {
				t.Errorf("expected IsBot(%q) == false, got true", ua)
			}
		}
	})

	t.Run("still flags real bots", func(t *testing.T) {
		realBots := []string{
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
			"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
			"some-crawler/1.0",
			"a-scraper-bot",
		}

		for _, ua := range realBots {
			match := DetectBot(ua)
			if match == nil {
				t.Errorf("expected %q to be detected as a bot, got nil", ua)
			}
			if !IsBot(ua) {
				t.Errorf("expected IsBot(%q) == true, got false", ua)
			}
			if !IsCrawler(ua) {
				t.Errorf("expected IsCrawler(%q) == true, got false", ua)
			}
		}

		// Genuinely suspicious bare scanner tokens from Generic Bot alternation
		for _, ua := range []string{"ZmEu", "Zeus"} {
			match := DetectBot(ua)
			if match == nil {
				t.Errorf("expected %q to be detected, got nil", ua)
			} else if match.Name != "Generic Bot" {
				t.Errorf("expected name 'Generic Bot' for %q, got %q", ua, match.Name)
			}
		}
	})

	t.Run("checks search bots specifically", func(t *testing.T) {
		ua := "Googlebot/2.1 (+http://www.google.com/bot.html)"
		if !IsSearchBot(ua) {
			t.Errorf("expected Googlebot to be identified as search bot")
		}
	})
}

func BenchmarkDetectBot_Browser(b *testing.B) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	b.ResetTimer()
	for b.Loop() {
		_ = DetectBot(ua)
	}
}

func BenchmarkDetectBot_Googlebot(b *testing.B) {
	ua := "Googlebot/2.1 (+http://www.google.com/bot.html)"
	b.ResetTimer()
	for b.Loop() {
		_ = DetectBot(ua)
	}
}

func BenchmarkDetectBot_IncludesBot(b *testing.B) {
	ua := "WireReaderBot/1.0"
	b.ResetTimer()
	for b.Loop() {
		_ = DetectBot(ua)
	}
}

