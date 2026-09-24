package uaparser

import (
	"regexp"
	"strings"
)

// Pre-compiled regex patterns for AI Crawlers, AI Assistants, and General Bots.
var (
	// AI Assistants represent user-initiated AI browsing sessions (e.g. browsing the web via ChatGPT or Claude).
	aiAssistantsRegex = regexp.MustCompile(`(?i)(ChatGPT-User|Claude-Web|Perplexity-User|MistralAI-User|Cohere-AI|DuckAssistBot|Gemini-Deep-Research|Amazon-Nova-Act)`)

	// AI Crawlers represent automated data collection agents used for LLM training and RAG ingestion.
	aiCrawlersRegex = regexp.MustCompile(`(?i)(GPTBot|OAI-SearchBot|ClaudeBot|Claude-SearchBot|anthropic-ai|PerplexityBot|Bytespider|TikTokSpider|Google-Extended|Google-NotebookLM|CloudVertexBot|FacebookBot|Meta-ExternalAgent|Meta-ExternalFetcher|DeepSeekBot|CCBot|Diffbot|Applebot-Extended|cohere-training-data-crawler|DataForSeoBot|FirecrawlAgent|KimiBot|v0Bot|xai-bot|YouBot|ChatGLM-Spider|HuggingFaceBot|TogetherBot|ReplicateBot|PetalBot|PanguBot|CoveoBot|Amazonbot|AI2Bot|Timpibot|Omgilibot)`)

	// General search crawlers and social media preview bots.
	searchAndSocialBotsRegex = regexp.MustCompile(`(?i)(Googlebot|bingbot|Baiduspider|YandexBot|DuckDuckBot|Slurp|Sogou|Exabot|Facebot|facebookexternalhit|Twitterbot|LinkedInBot|Slackbot|TelegramBot|WhatsApp|Discordbot|Pinterestbot|Applebot)`)

	// Automated CLI tools, libraries, and HTTP clients.
	cliAndLibraryRegex = regexp.MustCompile(`(?i)(curl/|Wget/|python-requests|aiohttp|urllib|Go-http-client|node-fetch|axios/|PostmanRuntime|insomnia/|Apache-HttpClient|Java/|Ruby|PHP/|Scrapy|HeadlessChrome|PhantomJS)`)
)

// Known bot token substrings for high-speed early detection.
var botSubstrings = []string{
	"bot", "spider", "crawl", "slurp", "fetch", "archive", "scraper", "headless",
}

// IsAICrawler reports whether the User-Agent belongs to an AI data crawler or search bot (e.g. GPTBot, ClaudeBot).
func IsAICrawler(ua string) bool {
	return aiCrawlersRegex.MatchString(ua)
}

// IsAIAssistant reports whether the User-Agent belongs to a user-facing AI assistant (e.g. ChatGPT-User, Claude-Web).
func IsAIAssistant(ua string) bool {
	return aiAssistantsRegex.MatchString(ua)
}

// IsBot reports whether the User-Agent belongs to any known bot, crawler, search engine, scraper, CLI tool, or AI agent.
func IsBot(ua string) bool {
	if ua == "" {
		return false
	}

	// Fast path check using lowercase tokens
	lower := strings.ToLower(ua)
	for _, sub := range botSubstrings {
		if strings.Contains(lower, sub) {
			return true
		}
	}

	// Regex check for AI agents, crawlers, CLI tools, and libraries
	return aiCrawlersRegex.MatchString(ua) ||
		aiAssistantsRegex.MatchString(ua) ||
		searchAndSocialBotsRegex.MatchString(ua) ||
		cliAndLibraryRegex.MatchString(ua)
}
