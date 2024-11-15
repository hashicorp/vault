package models
// Allows IT admind to set a predefined default search engine for MDM-Controlled devices
type WindowsSpotlightEnablementSettings int

const (
    // Spotlight on lock screen is not configured
    NOTCONFIGURED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS WindowsSpotlightEnablementSettings = iota
    // Disable Windows Spotlight on lock screen
    DISABLED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
    // Enable Windows Spotlight on lock screen
    ENABLED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
)

func (i WindowsSpotlightEnablementSettings) String() string {
    return []string{"notConfigured", "disabled", "enabled"}[i]
}
func ParseWindowsSpotlightEnablementSettings(v string) (any, error) {
    result := NOTCONFIGURED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
        case "disabled":
            result = DISABLED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
        case "enabled":
            result = ENABLED_WINDOWSSPOTLIGHTENABLEMENTSETTINGS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsSpotlightEnablementSettings(values []WindowsSpotlightEnablementSettings) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsSpotlightEnablementSettings) isMultiValue() bool {
    return false
}
