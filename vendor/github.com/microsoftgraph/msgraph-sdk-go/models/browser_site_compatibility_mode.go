package models
type BrowserSiteCompatibilityMode int

const (
    // Loads the site using default compatibility mode.
    DEFAULT_BROWSERSITECOMPATIBILITYMODE BrowserSiteCompatibilityMode = iota
    // Loads the site in internetExplorer8 Enterprise Mode
    INTERNETEXPLORER8ENTERPRISE_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer7 Enterprise Mode
    INTERNETEXPLORER7ENTERPRISE_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer11
    INTERNETEXPLORER11_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer10
    INTERNETEXPLORER10_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer9
    INTERNETEXPLORER9_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer8
    INTERNETEXPLORER8_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer7
    INTERNETEXPLORER7_BROWSERSITECOMPATIBILITYMODE
    // Loads the site in internetExplorer5
    INTERNETEXPLORER5_BROWSERSITECOMPATIBILITYMODE
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shouldn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSITECOMPATIBILITYMODE
)

func (i BrowserSiteCompatibilityMode) String() string {
    return []string{"default", "internetExplorer8Enterprise", "internetExplorer7Enterprise", "internetExplorer11", "internetExplorer10", "internetExplorer9", "internetExplorer8", "internetExplorer7", "internetExplorer5", "unknownFutureValue"}[i]
}
func ParseBrowserSiteCompatibilityMode(v string) (any, error) {
    result := DEFAULT_BROWSERSITECOMPATIBILITYMODE
    switch v {
        case "default":
            result = DEFAULT_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer8Enterprise":
            result = INTERNETEXPLORER8ENTERPRISE_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer7Enterprise":
            result = INTERNETEXPLORER7ENTERPRISE_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer11":
            result = INTERNETEXPLORER11_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer10":
            result = INTERNETEXPLORER10_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer9":
            result = INTERNETEXPLORER9_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer8":
            result = INTERNETEXPLORER8_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer7":
            result = INTERNETEXPLORER7_BROWSERSITECOMPATIBILITYMODE
        case "internetExplorer5":
            result = INTERNETEXPLORER5_BROWSERSITECOMPATIBILITYMODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSITECOMPATIBILITYMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSiteCompatibilityMode(values []BrowserSiteCompatibilityMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSiteCompatibilityMode) isMultiValue() bool {
    return false
}
