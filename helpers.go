package uaparser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	appModelRegex        = regexp.MustCompile(`(?i)Model=([^;)]+)`)
	appManufacturerRegex = regexp.MustCompile(`(?i)Manufacturer=([^;)]+)`)
	singleNameVersionReg = regexp.MustCompile(`^[^\/]+\/[\d.]+$`)
	cleanVersionRegex    = regexp.MustCompile(`[^\d\.]`)

	// Brand patterns with precedence
	brandPatterns = []struct {
		regex *regexp.Regexp
		brand string
	}{
		{regexp.MustCompile(`(?i)\bxiaomi\b`), "Xiaomi"},
		{regexp.MustCompile(`(?i)\bredmi\b`), "Xiaomi"},
		{regexp.MustCompile(`(?i)\bpoco\b`), "Xiaomi"},
		{regexp.MustCompile(`(?i)\bsamsung\b`), "Samsung"},
		{regexp.MustCompile(`(?i)\bgalaxy\b`), "Samsung"},
		{regexp.MustCompile(`(?i)\bhuawei\b`), "Huawei"},
		{regexp.MustCompile(`(?i)\bhonor\b`), "Honor"},
		{regexp.MustCompile(`(?i)\boppo\b`), "OPPO"},
		{regexp.MustCompile(`(?i)\bvivo\b`), "Vivo"},
		{regexp.MustCompile(`(?i)\biqoo\b`), "Vivo"},
		{regexp.MustCompile(`(?i)\boneplus\b`), "OnePlus"},
		{regexp.MustCompile(`(?i)\bgoogle\b`), "Google"},
		{regexp.MustCompile(`(?i)\bpixel\b`), "Google"},
		{regexp.MustCompile(`(?i)\brealme\b`), "Realme"},
		{regexp.MustCompile(`(?i)\bmotorola\b`), "Motorola"},
		{regexp.MustCompile(`(?i)\bmoto\s`), "Motorola"},
		{regexp.MustCompile(`(?i)\bnokia\b`), "Nokia"},
		{regexp.MustCompile(`(?i)\bsony\b`), "Sony"},
		{regexp.MustCompile(`(?i)\bxperia\b`), "Sony"},
		{regexp.MustCompile(`(?i)\bnothing\b`), "Nothing"},
		{regexp.MustCompile(`(?i)\bapple\b`), "Apple"},
		{regexp.MustCompile(`(?i)\biphone\b`), "Apple"},
		{regexp.MustCompile(`(?i)\bipad\b`), "Apple"},
		{regexp.MustCompile(`(?i)\blg[- /]`), "LG"},
		{regexp.MustCompile(`(?i)\bzte\b`), "ZTE"},
		{regexp.MustCompile(`(?i)\blenovo\b`), "Lenovo"},
		{regexp.MustCompile(`(?i)\basus\b`), "ASUS"},
		{regexp.MustCompile(`(?i)\btcl\b`), "TCL"},
	}
)

// Majorize extracts the major version number from a version string (e.g. "124.0.6367.60" -> "124").
func Majorize(version string) string {
	if version == "" {
		return ""
	}
	clean := cleanVersionRegex.ReplaceAllString(version, "")
	parts := strings.Split(clean, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// AppStyleInfo contains device info parsed from app-style UA strings (e.g., "Model=Redmi Note 8 Pro; Manufacturer=Xiaomi").
type AppStyleInfo struct {
	Model        string
	Manufacturer string
}

// ExtractAppStyleInfo extracts device and model metadata from custom mobile application User-Agents.
func ExtractAppStyleInfo(ua string) AppStyleInfo {
	var info AppStyleInfo
	if m := appModelRegex.FindStringSubmatch(ua); len(m) > 1 {
		info.Model = strings.TrimSpace(m[1])
	}
	if m := appManufacturerRegex.FindStringSubmatch(ua); len(m) > 1 {
		info.Manufacturer = strings.TrimSpace(m[1])
	}
	return info
}

// DetectBrand determines the device brand from a User-Agent or model name.
func DetectBrand(ua, model string) string {
	combined := ua + " " + model
	for _, entry := range brandPatterns {
		if entry.regex.MatchString(combined) {
			return entry.brand
		}
	}
	return ""
}

// GetOutlookEdition translates raw Microsoft Outlook version numbers into developer-friendly edition titles.
func GetOutlookEdition(name, version string) string {
	if name == "" || version == "" {
		return name
	}
	cleanName := strings.ToLower(name)
	cleanName = strings.ReplaceAll(cleanName, "microsoft ", "")

	// 1. Handle Mac Separately
	if cleanName == "macoutlook" {
		major, _ := strconv.Atoi(Majorize(version))
		if major >= 16 {
			return "Outlook for Mac (Modern)"
		}
		return "Outlook for Mac (Legacy)"
	}

	// 2. Handle Windows Outlook
	if cleanName == "outlook" {
		parts := strings.Split(version, ".")
		major, _ := strconv.Atoi(parts[0])
		var build int
		if len(parts) >= 3 {
			build, _ = strconv.Atoi(parts[2])
		}

		if major == 15 {
			return "Outlook 2013"
		}
		if major == 14 {
			return "Outlook 2010"
		}
		if major == 12 {
			return "Outlook 2007"
		}
		if major < 12 && major > 0 {
			return "Outlook (Legacy)"
		}
		if major == 16 {
			if build > 0 && build < 10000 {
				return "Outlook 2016 (MSI / Volume License)"
			}
			return "Outlook 365 / 2019+ (Modern)"
		}
	}

	return name
}

// IsServerUA checks if a User-Agent represents a server-to-server HTTP client or headless tool.
func IsServerUA(ua string, res *Result) bool {
	if singleNameVersionReg.MatchString(ua) {
		return true
	}
	if res.OS.Name == "" && res.Browser.Name == "" && res.Device.Vendor == "" && res.Device.Model == "" {
		return true
	}
	return false
}
