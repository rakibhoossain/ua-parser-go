package uaparser

import (
	"regexp"
	"strings"
)

type browserRule struct {
	regex   *regexp.Regexp
	handler func(matches []string) Browser
}

type osRule struct {
	regex   *regexp.Regexp
	handler func(matches []string) OS
}

type deviceRule struct {
	regex   *regexp.Regexp
	handler func(matches []string) Device
}

type engineRule struct {
	regex   *regexp.Regexp
	handler func(matches []string) Engine
}

type cpuRule struct {
	regex   *regexp.Regexp
	handler func(matches []string) CPU
}

var (
	browserRules []browserRule
	osRules      []osRule
	deviceRules  []deviceRule
	engineRules  []engineRule
	cpuRules     []cpuRule
)

func init() {
	initBrowserRules()
	initOSRules()
	initDeviceRules()
	initEngineRules()
	initCPURules()
}

func initBrowserRules() {
	browserRules = []browserRule{
		// Chrome for iOS / Android
		{
			regex: regexp.MustCompile(`(?i)\b(?:crmo|crios)/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Mobile Chrome", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Edge WebView
		{
			regex: regexp.MustCompile(`(?i)webview.+edge/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Edge WebView", Version: m[1], Major: Majorize(m[1]), Type: BrowserTypeInApp}
			},
		},
		// Edge (Chromium / iOS / Android)
		{
			regex: regexp.MustCompile(`(?i)edg(?:e|ios|a)?/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Edge", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Opera Mini
		{
			regex: regexp.MustCompile(`(?i)(opera mini)/([-\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Opera Mini", Version: m[2], Major: Majorize(m[2])}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)opios[/ ]+([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Opera Mini", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Opera GX
		{
			regex: regexp.MustCompile(`(?i)\bop(?:rg)?x/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Opera GX", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Opera
		{
			regex: regexp.MustCompile(`(?i)\bopr/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Opera", Version: m[1], Major: Majorize(m[1])}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)(opera)(?:.+version/|[/ ]+)([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Opera", Version: m[2], Major: Majorize(m[2])}
			},
		},
		// Samsung Internet
		{
			regex: regexp.MustCompile(`(?i)samsungbrowser/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Samsung Internet", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Baidu
		{
			regex: regexp.MustCompile(`(?i)\bb[ai]*d(?:uhd|[ub]*[aekoprswx]{5,6})[/ ]?([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Baidu", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Brave
		{
			regex: regexp.MustCompile(`(?i)(brave)(?: chrome)?/([\d\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Brave", Version: m[2], Major: Majorize(m[2])}
			},
		},
		// Vivaldi
		{
			regex: regexp.MustCompile(`(?i)vivaldi/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Vivaldi", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// DuckDuckGo
		{
			regex: regexp.MustCompile(`(?i)(?:duckduckgo|ddg)/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "DuckDuckGo", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// UCBrowser
		{
			regex: regexp.MustCompile(`(?i)(?:\buc? ?browser|(?:juc.+)ucweb| ucpc)[/ ]?([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "UCBrowser", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// WeChat
		{
			regex: regexp.MustCompile(`(?i)micromessenger/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "WeChat", Version: m[1], Major: Majorize(m[1]), Type: BrowserTypeInApp}
			},
		},
		// QQBrowser
		{
			regex: regexp.MustCompile(`(?i)m?qqbrowser(?:lite)?/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "QQBrowser", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Yandex
		{
			regex: regexp.MustCompile(`(?i)yabrowser/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Yandex", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Huawei Browser
		{
			regex: regexp.MustCompile(`(?i)huaweibrowser/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Huawei Browser", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// MIUI Browser
		{
			regex: regexp.MustCompile(`(?i)miuibrowser/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "MIUI Browser", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Outlook
		{
			regex: regexp.MustCompile(`(?i)(microsoft |mac)?outlook[/ ]?([\d\.]+)`),
			handler: func(m []string) Browser {
				name := GetOutlookEdition("Outlook", m[2])
				return Browser{Name: name, Version: m[2], Major: Majorize(m[2]), Type: BrowserTypeEmail}
			},
		},
		// Facebook
		{
			regex: regexp.MustCompile(`(?i)(?:fbav|fban)/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Facebook", Version: m[1], Major: Majorize(m[1]), Type: BrowserTypeInApp}
			},
		},
		// Instagram
		{
			regex: regexp.MustCompile(`(?i)instagram[/ ]?([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Instagram", Version: m[1], Major: Majorize(m[1]), Type: BrowserTypeInApp}
			},
		},
		// TikTok
		{
			regex: regexp.MustCompile(`(?i)(?:app_type/)?(musical_ly|tiktok)[/ ]?([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "TikTok", Version: m[2], Major: Majorize(m[2]), Type: BrowserTypeInApp}
			},
		},
		// Headless Chrome
		{
			regex: regexp.MustCompile(`(?i)headlesschrome(?:/([\w\.]+)| )`),
			handler: func(m []string) Browser {
				v := ""
				if len(m) > 1 {
					v = m[1]
				}
				return Browser{Name: "Chrome Headless", Version: v, Major: Majorize(v), Type: BrowserTypeCrawler}
			},
		},
		// Chrome WebView
		{
			regex: regexp.MustCompile(`(?i)\bwv\b.+chrome/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Chrome WebView", Version: m[1], Major: Majorize(m[1]), Type: BrowserTypeInApp}
			},
		},
		// Chrome
		{
			regex: regexp.MustCompile(`(?i)(?:chrome|crios)/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Chrome", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Chromium
		{
			regex: regexp.MustCompile(`(?i)chromium/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Chromium", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Firefox Focus
		{
			regex: regexp.MustCompile(`(?i)focus/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Firefox Focus", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Firefox
		{
			regex: regexp.MustCompile(`(?i)(?:firefox|fxios)/([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "Firefox", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Safari Mobile
		{
			regex: regexp.MustCompile(`(?i)version/([\w\.]+).+mobile/\w+ safari`),
			handler: func(m []string) Browser {
				return Browser{Name: "Mobile Safari", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Safari
		{
			regex: regexp.MustCompile(`(?i)version/([\w\.]+).+safari`),
			handler: func(m []string) Browser {
				return Browser{Name: "Safari", Version: m[1], Major: Majorize(m[1])}
			},
		},
		// Internet Explorer
		{
			regex: regexp.MustCompile(`(?i)(?:msie |trident.+rv:)([\w\.]+)`),
			handler: func(m []string) Browser {
				return Browser{Name: "IE", Version: m[1], Major: Majorize(m[1])}
			},
		},
	}
}

func initOSRules() {
	osRules = []osRule{
		// Windows
		{
			regex: regexp.MustCompile(`(?i)windows nt 10\.0`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "10"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 6\.3`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "8.1"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 6\.2`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "8"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 6\.1`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "7"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 6\.0`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "Vista"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 5\.1`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "XP"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt 5\.0`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "2000"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)windows nt ([\d\.]+)`),
			handler: func(m []string) OS {
				return OS{Name: "Windows", Version: "NT " + m[1]}
			},
		},
		// iOS / iPadOS
		{
			regex: regexp.MustCompile(`(?i)ip[honead]+;.+os ([\d_]+)`),
			handler: func(m []string) OS {
				v := strings.ReplaceAll(m[1], "_", ".")
				return OS{Name: "iOS", Version: v}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)ipados[ /]([\d_]+)`),
			handler: func(m []string) OS {
				v := strings.ReplaceAll(m[1], "_", ".")
				return OS{Name: "iOS", Version: v}
			},
		},
		// macOS
		{
			regex: regexp.MustCompile(`(?i)mac os x ([\d_\.]+)`),
			handler: func(m []string) OS {
				v := strings.ReplaceAll(m[1], "_", ".")
				return OS{Name: "macOS", Version: v}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)macintosh;.+mac os x`),
			handler: func(m []string) OS {
				return OS{Name: "macOS"}
			},
		},
		// Android
		{
			regex: regexp.MustCompile(`(?i)android[ /]?([\d\.]*)`),
			handler: func(m []string) OS {
				v := m[1]
				if v == "" {
					v = ""
				}
				return OS{Name: "Android", Version: v}
			},
		},
		// Chrome OS
		{
			regex: regexp.MustCompile(`(?i)cros [^\s]+ ([\d\.]+)`),
			handler: func(m []string) OS {
				return OS{Name: "Chrome OS", Version: m[1]}
			},
		},
		// Linux Distros
		{
			regex: regexp.MustCompile(`(?i)(ubuntu|debian|fedora|centos|gentoo|arch|slackware|redhat|suse|mint|freebsd|openbsd|netbsd)`),
			handler: func(m []string) OS {
				name := strings.Title(strings.ToLower(m[1]))
				return OS{Name: name}
			},
		},
		// Generic Linux
		{
			regex: regexp.MustCompile(`(?i)linux`),
			handler: func(m []string) OS {
				return OS{Name: "Linux"}
			},
		},
		// SmartTV / Others
		{
			regex: regexp.MustCompile(`(?i)tizen[ /]?([\d\.]*)`),
			handler: func(m []string) OS {
				return OS{Name: "Tizen", Version: m[1]}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)webos[ /]?([\d\.]*)`),
			handler: func(m []string) OS {
				return OS{Name: "WebOS", Version: m[1]}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)playstation ([\w\d]+)`),
			handler: func(m []string) OS {
				return OS{Name: "PlayStation", Version: m[1]}
			},
		},
	}
}

func initDeviceRules() {
	deviceRules = []deviceRule{
		// Apple iPad
		{
			regex: regexp.MustCompile(`(?i)\b(ipad)[\d,]*`),
			handler: func(m []string) Device {
				return Device{Vendor: "Apple", Model: "iPad", Type: DeviceTablet}
			},
		},
		// Apple iPhone / iPod
		{
			regex: regexp.MustCompile(`(?i)\b(iphone|ipod)\b`),
			handler: func(m []string) Device {
				model := "iPhone"
				if strings.EqualFold(m[1], "ipod") {
					model = "iPod"
				}
				return Device{Vendor: "Apple", Model: model, Type: DeviceMobile}
			},
		},
		// Apple Macintosh
		{
			regex: regexp.MustCompile(`(?i)\bmacintosh\b`),
			handler: func(m []string) Device {
				return Device{Vendor: "Apple", Model: "Macintosh", Type: DeviceDesktop}
			},
		},
		// Samsung Tablets (SM-T, GT-P, Galaxy Tab)
		{
			regex: regexp.MustCompile(`(?i)\b(sm-t[0-9]+|gt-p[0-9]+|galaxy tab[\w ]*)\b`),
			handler: func(m []string) Device {
				return Device{Vendor: "Samsung", Model: m[1], Type: DeviceTablet}
			},
		},
		// Samsung Mobile (SM-[A-Z0-9]+, Galaxy S/A/Note)
		{
			regex: regexp.MustCompile(`(?i)\b((?:sm-|gt-)[a-z0-9]+)\b`),
			handler: func(m []string) Device {
				lower := strings.ToLower(m[1])
				devType := DeviceMobile
				if strings.HasPrefix(lower, "sm-t") || strings.HasPrefix(lower, "gt-p") {
					devType = DeviceTablet
				} else if strings.HasPrefix(lower, "sm-l") || strings.HasPrefix(lower, "sm-r") {
					devType = DeviceWearable
				}
				return Device{Vendor: "Samsung", Model: m[1], Type: devType}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)\b(galaxy\s*(?:s|a|m|note|z\s*(?:fold|flip))[0-9]+[a-z0-9]*)\b`),
			handler: func(m []string) Device {
				return Device{Vendor: "Samsung", Model: m[1], Type: DeviceMobile}
			},
		},
		// LG Mobile
		{
			regex: regexp.MustCompile(`(?i)\b(lg-[a-z0-9]+)\b`),
			handler: func(m []string) Device {
				return Device{Vendor: "LG", Model: m[1], Type: DeviceMobile}
			},
		},
		// Google Pixel
		{
			regex: regexp.MustCompile(`(?i)\b(pixel(?:\s+[0-9a-z]+)?)(?:\s+build|[;\)])`),
			handler: func(m []string) Device {
				return Device{Vendor: "Google", Model: strings.TrimSpace(m[1]), Type: DeviceMobile}
			},
		},
		// Xiaomi / Redmi / POCO
		{
			regex: regexp.MustCompile(`(?i)\b(redmi\s*(?:note|k|[0-9])[\w ]*|poco\s*[a-z0-9 ]+|mi\s*[0-9][\w ]*)\b`),
			handler: func(m []string) Device {
				return Device{Vendor: "Xiaomi", Model: strings.TrimSpace(m[1]), Type: DeviceMobile}
			},
		},
		// Huawei / Honor
		{
			regex: regexp.MustCompile(`(?i)\b(huawei\s*[a-z0-9\- ]+|honor\s*[a-z0-9\- ]+)\b`),
			handler: func(m []string) Device {
				vendor := "Huawei"
				if strings.Contains(strings.ToLower(m[1]), "honor") {
					vendor = "Honor"
				}
				return Device{Vendor: vendor, Model: strings.TrimSpace(m[1]), Type: DeviceMobile}
			},
		},
		// SmartTV devices
		{
			regex: regexp.MustCompile(`(?i)(smart-tv|hbbtv|appletv|crkey|googletv|roku)`),
			handler: func(m []string) Device {
				vendor := ""
				match := strings.ToLower(m[1])
				if strings.Contains(match, "appletv") {
					vendor = "Apple"
				} else if strings.Contains(match, "crkey") || strings.Contains(match, "googletv") {
					vendor = "Google"
				} else if strings.Contains(match, "roku") {
					vendor = "Roku"
				}
				return Device{Vendor: vendor, Model: m[1], Type: DeviceSmartTV}
			},
		},
		// Consoles
		{
			regex: regexp.MustCompile(`(?i)(playstation|xbox|nintendo switch)`),
			handler: func(m []string) Device {
				vendor := ""
				match := strings.ToLower(m[1])
				if strings.Contains(match, "playstation") {
					vendor = "Sony"
				} else if strings.Contains(match, "xbox") {
					vendor = "Microsoft"
				} else if strings.Contains(match, "nintendo") {
					vendor = "Nintendo"
				}
				return Device{Vendor: vendor, Model: m[1], Type: DeviceConsole}
			},
		},
	}
}

func initEngineRules() {
	engineRules = []engineRule{
		{
			regex: regexp.MustCompile(`(?i)(blink|chrome|crios|crmo)`),
			handler: func(m []string) Engine {
				return Engine{Name: "Blink"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)applewebkit/([\d\.]+)`),
			handler: func(m []string) Engine {
				return Engine{Name: "WebKit", Version: m[1]}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)rv:([\d\.]+).+gecko/`),
			handler: func(m []string) Engine {
				return Engine{Name: "Gecko", Version: m[1]}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)trident/([\d\.]+)`),
			handler: func(m []string) Engine {
				return Engine{Name: "Trident", Version: m[1]}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)presto/([\d\.]+)`),
			handler: func(m []string) Engine {
				return Engine{Name: "Presto", Version: m[1]}
			},
		},
	}
}

func initCPURules() {
	cpuRules = []cpuRule{
		{
			regex: regexp.MustCompile(`(?i)\b(arm64|aarch64)\b`),
			handler: func(m []string) CPU {
				return CPU{Architecture: "arm64"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)\b(armv\d+[a-z]*|arm)\b`),
			handler: func(m []string) CPU {
				return CPU{Architecture: "arm"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)\b(x86_64|amd64|win64|x64|wow64)\b`),
			handler: func(m []string) CPU {
				return CPU{Architecture: "amd64"}
			},
		},
		{
			regex: regexp.MustCompile(`(?i)\b(i[3-6]86|x86|ia32)\b`),
			handler: func(m []string) CPU {
				return CPU{Architecture: "ia32"}
			},
		},
	}
}

func matchBrowser(ua string) Browser {
	for _, rule := range browserRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			return rule.handler(matches)
		}
	}
	return Browser{}
}

func matchOS(ua string) OS {
	for _, rule := range osRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			return rule.handler(matches)
		}
	}
	return OS{}
}

func matchDevice(ua string) Device {
	for _, rule := range deviceRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			return rule.handler(matches)
		}
	}
	return Device{}
}

func matchEngine(ua string) Engine {
	for _, rule := range engineRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			return rule.handler(matches)
		}
	}
	return Engine{}
}

func matchCPU(ua string) CPU {
	for _, rule := range cpuRules {
		if matches := rule.regex.FindStringSubmatch(ua); len(matches) > 0 {
			return rule.handler(matches)
		}
	}
	return CPU{}
}
