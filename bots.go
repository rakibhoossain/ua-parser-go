package uaparser

import (
	"strings"

	"github.com/rakibhoossain/ua-parser-go/bots"
)

// DetectBot identifies the User-Agent using the authoritative bot database,
// returning detailed bot metadata (Name, Category, URL, Producer) or nil if not a bot.
func DetectBot(ua string) *bots.BotMatch {
	return bots.DetectBot(ua)
}

// IsBot reports whether the User-Agent belongs to any known bot, crawler, search engine, or automated scraper.
func IsBot(ua string) bool {
	return bots.IsBot(ua)
}

// IsCrawler reports whether the User-Agent represents a known search spider, web crawler, or SEO bot
// (Googlebot, Bingbot, BingPreview, Yandex, Baidu, Screaming Frog, FacebookExternalHit, GPTBot, ClaudeBot, etc.).
func IsCrawler(ua string) bool {
	return bots.IsCrawler(ua)
}

// IsSearchBot reports whether the User-Agent is explicitly a search engine bot.
func IsSearchBot(ua string) bool {
	return bots.IsSearchBot(ua)
}

// IsAICrawler reports whether the User-Agent belongs to an AI data scraper or crawler (e.g. GPTBot, ClaudeBot, Bytespider).
func IsAICrawler(ua string) bool {
	b := bots.DetectBot(ua)
	if b == nil {
		return false
	}
	cat := strings.ToLower(b.Category)
	return strings.Contains(cat, "ai") &&
		(strings.Contains(cat, "scraper") || strings.Contains(cat, "crawler") || strings.Contains(cat, "data"))
}

// IsAIAssistant reports whether the User-Agent belongs to a user-facing AI assistant (e.g. ChatGPT-User, Claude-Web).
func IsAIAssistant(ua string) bool {
	b := bots.DetectBot(ua)
	if b == nil {
		return false
	}
	return strings.EqualFold(b.Category, "ai assistant")
}
