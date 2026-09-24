package uaparser

// DeviceType represents the general category of the client hardware.
type DeviceType string

const (
	DeviceDesktop  DeviceType = "desktop"
	DeviceMobile   DeviceType = "mobile"
	DeviceTablet   DeviceType = "tablet"
	DeviceSmartTV  DeviceType = "smarttv"
	DeviceConsole  DeviceType = "console"
	DeviceWearable DeviceType = "wearable"
	DeviceXR       DeviceType = "xr"
	DeviceEmbedded DeviceType = "embedded"
)

// BrowserType represents the category of the browser/client application.
type BrowserType string

const (
	BrowserTypeRegular     BrowserType = ""
	BrowserTypeInApp       BrowserType = "inapp"
	BrowserTypeCrawler     BrowserType = "crawler"
	BrowserTypeFetcher     BrowserType = "fetcher"
	BrowserTypeCLI         BrowserType = "cli"
	BrowserTypeEmail       BrowserType = "email"
	BrowserTypeLibrary     BrowserType = "library"
	BrowserTypeMediaPlayer BrowserType = "mediaplayer"
)

// Browser contains parsed browser metadata.
type Browser struct {
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Major   string      `json:"major"`
	Type    BrowserType `json:"type,omitempty"`
}

// OS contains parsed operating system metadata.
type OS struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Device contains parsed hardware device metadata.
type Device struct {
	Vendor string     `json:"vendor"`
	Model  string     `json:"model"`
	Type   DeviceType `json:"type"`
}

// Engine contains parsed browser rendering engine metadata.
type Engine struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CPU contains parsed CPU architecture metadata.
type CPU struct {
	Architecture string `json:"architecture"`
}

// ClientHints represents structured User-Agent Client Hints (UACH) values
// parsed from HTTP headers (e.g. Sec-CH-UA, Sec-CH-UA-Platform, etc.).
type ClientHints struct {
	Brands          []Brand `json:"brands,omitempty"`
	FullVersionList []Brand `json:"full_version_list,omitempty"`
	Mobile          *bool   `json:"mobile,omitempty"`
	Model           string  `json:"model,omitempty"`
	Platform        string  `json:"platform,omitempty"`
	PlatformVersion string  `json:"platform_version,omitempty"`
	Architecture    string  `json:"architecture,omitempty"`
	Bitness         string  `json:"bitness,omitempty"`
	FormFactors     []string`json:"form_factors,omitempty"`
}

// Brand represents a browser brand token and version from Sec-CH-UA headers.
type Brand struct {
	Brand   string `json:"brand"`
	Version string `json:"version"`
}

// Result is the consolidated parsed output representing all detected details.
type Result struct {
	UA          string      `json:"ua"`
	Browser     Browser     `json:"browser"`
	OS          OS          `json:"os"`
	Device      Device      `json:"device"`
	Engine      Engine      `json:"engine"`
	CPU         CPU         `json:"cpu"`
	ClientHints ClientHints `json:"client_hints,omitempty"`
	IsBot       bool        `json:"is_bot"`
	IsServer    bool        `json:"is_server"`
	IsFrozen    bool        `json:"is_frozen"`
}
