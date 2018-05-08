package ldaputil

type ConfigEntry struct {
	Url           string `json:"url"`
	UserDN        string `json:"userdn"`
	GroupDN       string `json:"groupdn"`
	GroupFilter   string `json:"groupfilter"`
	GroupAttr     string `json:"groupattr"`
	UPNDomain     string `json:"upndomain"`
	UserAttr      string `json:"userattr"`
	Certificate   string `json:"certificate"`
	InsecureTLS   bool   `json:"insecure_tls"`
	StartTLS      bool   `json:"starttls"`
	BindDN        string `json:"binddn"`
	BindPassword  string `json:"bindpass"`
	DenyNullBind  bool   `json:"deny_null_bind"`
	DiscoverDN    bool   `json:"discoverdn"`
	TLSMinVersion string `json:"tls_min_version"`
	TLSMaxVersion string `json:"tls_max_version"`

	// This json tag deviates from snake case because there was a past issue
	// where the tag was being ignored, causing it to be jsonified as "CaseSensitiveNames".
	// To continue reading in users' previously stored values,
	// we chose to carry that forward.
	CaseSensitiveNames *bool `json:"CaseSensitiveNames,omitempty"`
}
