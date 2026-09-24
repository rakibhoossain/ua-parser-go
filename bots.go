package uaparser

import (
	"regexp"

	"github.com/rakibhoossain/ua-parser-go/bots"
)

// Pre-compiled regex patterns for AI Crawlers, AI Assistants, and General Bots.
var (
	// AI Assistants represent user-initiated AI browsing sessions (e.g. browsing the web via ChatGPT or Claude).
	aiAssistantsRegex = regexp.MustCompile(`(?i)(ChatGPT-User|Claude-Web|Perplexity-User|MistralAI-User|Cohere-AI|DuckAssistBot|Gemini-Deep-Research|Amazon-Nova-Act)`)

	// AI Crawlers represent automated data collection agents used for LLM training and RAG ingestion.
	aiCrawlersRegex = regexp.MustCompile(`(?i)(GPTBot|OAI-SearchBot|ClaudeBot|Claude-SearchBot|anthropic-ai|PerplexityBot|Bytespider|TikTokSpider|Google-Extended|Google-NotebookLM|CloudVertexBot|FacebookBot|Meta-ExternalAgent|Meta-ExternalFetcher|DeepSeekBot|CCBot|Diffbot|Applebot-Extended|cohere-training-data-crawler|DataForSeoBot|FirecrawlAgent|KimiBot|v0Bot|xai-bot|YouBot|ChatGLM-Spider|HuggingFaceBot|TogetherBot|ReplicateBot|PetalBot|PanguBot|CoveoBot|Amazonbot|AI2Bot|Timpibot|Omgilibot)`)

	// General search crawlers and social media preview bots.
	searchAndSocialBotsRegex = regexp.MustCompile(`(?i)(Googlebot|bingbot|bingpreview|Baiduspider|YandexBot|DuckDuckBot|Slurp|Sogou|Exabot|Facebot|facebookexternalhit|Twitterbot|LinkedInBot|Slackbot|TelegramBot|WhatsApp|Discordbot|Pinterest|Applebot|SemrushBot|AhrefsBot|MJ12bot|DotBot|screaming frog|seznambot|archive\.org_bot|ia_archiver)`)

	// Automated CLI tools, libraries, and HTTP clients.
	cliAndLibraryRegex = regexp.MustCompile(`(?i)(curl/|Wget/|python-requests|aiohttp|urllib|Go-http-client|node-fetch|axios/|PostmanRuntime|insomnia/|Apache-HttpClient|Java/|Ruby|PHP/|Scrapy|HeadlessChrome|PhantomJS)`)
)

// DetectBot identifies the User-Agent using both high-speed token lookups and regex patterns,
// returning detailed bot metadata (name, category, producer, url) or nil if not a bot.
func DetectBot(ua string) *bots.BotMatch {
	if ua == "" {
		return nil
	}
	if b := bots.DetectBot(ua); b != nil {
		return b
	}
	if IsAICrawler(ua) {
		return &bots.BotMatch{Name: "AI Crawler", Category: "AI Crawler"}
	}
	if IsAIAssistant(ua) {
		return &bots.BotMatch{Name: "AI Assistant", Category: "AI Assistant"}
	}
	if cliAndLibraryRegex.MatchString(ua) {
		return &bots.BotMatch{Name: "HTTP Client / CLI", Category: "Library"}
	}
	return nil
}

// IsAICrawler reports whether the User-Agent belongs to an AI data crawler or search bot (e.g. GPTBot, ClaudeBot).
func IsAICrawler(ua string) bool {
	return aiCrawlersRegex.MatchString(ua)
}

// IsAIAssistant reports whether the User-Agent belongs to a user-facing AI assistant (e.g. ChatGPT-User, Claude-Web).
func IsAIAssistant(ua string) bool {
	return aiAssistantsRegex.MatchString(ua)
}

// IsCrawler reports whether the User-Agent represents a known search spider, web crawler, or SEO bot
// (Googlebot, Bingbot, BingPreview, Yandex, Baidu, Screaming Frog, FacebookExternalHit, GPTBot, ClaudeBot, etc.).
// This function is ideal for crawler-bypass hooks to prevent spiders from polluting analytics funnels.
func IsCrawler(ua string) bool {
	if ua == "" {
		return false
	}
	return bots.IsCrawler(ua) || searchAndSocialBotsRegex.MatchString(ua) || aiCrawlersRegex.MatchString(ua)
}

// IsSearchBot reports whether the User-Agent is explicitly a search engine bot.
func IsSearchBot(ua string) bool {
	if ua == "" {
		return false
	}
	return bots.IsSearchBot(ua) || searchAndSocialBotsRegex.MatchString(ua)
}

// IsBot reports whether the User-Agent belongs to any known bot, crawler, search engine, scraper, CLI tool, or AI agent.
func IsBot(ua string) bool {
	if ua == "" {
		return false
	}
	if bots.IsBot(ua) {
		return true
	}
	return aiCrawlersRegex.MatchString(ua) ||
		aiAssistantsRegex.MatchString(ua) ||
		searchAndSocialBotsRegex.MatchString(ua) ||
		cliAndLibraryRegex.MatchString(ua)
}
