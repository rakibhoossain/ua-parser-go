package uaparser_test

import (
	"fmt"
	"net/http"

	uaparser "github.com/rakibhoossain/ua-parser-go"
	"github.com/rakibhoossain/ua-parser-go/referrer"
)

func ExampleParse() {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	res := uaparser.Parse(ua)

	fmt.Printf("Browser: %s %s\n", res.Browser.Name, res.Browser.Major)
	fmt.Printf("OS: %s %s\n", res.OS.Name, res.OS.Version)
	fmt.Printf("Device: %s (%s)\n", res.Device.Vendor, res.Device.Type)

	// Output:
	// Browser: Chrome 124
	// OS: macOS 10.15.7
	// Device: Apple (desktop)
}

func ExampleParseRequest() {
	req, _ := http.NewRequest("GET", "https://analytics.example.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1")

	res := uaparser.ParseRequest(req)
	fmt.Printf("Browser: %s\n", res.Browser.Name)
	fmt.Printf("Device: %s %s\n", res.Device.Vendor, res.Device.Model)

	// Output:
	// Browser: Mobile Safari
	// Device: Apple iPhone
}

func Example_referrer() {
	refURL := "https://chatgpt.com/c/66f29ab0-1234"
	ref := referrer.Parse(refURL)

	fmt.Printf("Name: %s\n", ref.Name)
	fmt.Printf("Type: %s\n", ref.Type)
	fmt.Printf("Domain: %s\n", ref.Domain)
	fmt.Printf("Favicon: %s\n", ref.FaviconURL)

	// Output:
	// Name: ChatGPT
	// Type: ai
	// Domain: chatgpt.com
	// Favicon: https://icons.duckduckgo.com/ip3/chatgpt.com.ico
}
