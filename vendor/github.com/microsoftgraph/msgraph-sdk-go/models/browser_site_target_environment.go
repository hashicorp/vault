package models
type BrowserSiteTargetEnvironment int

const (
    // Open in Internet Explorer Mode
    INTERNETEXPLORERMODE_BROWSERSITETARGETENVIRONMENT BrowserSiteTargetEnvironment = iota
    // Open in standalone Internet Explorer 11
    INTERNETEXPLORER11_BROWSERSITETARGETENVIRONMENT
    // Open in Microsoft Edge
    MICROSOFTEDGE_BROWSERSITETARGETENVIRONMENT
    // Configurable type
    CONFIGURABLE_BROWSERSITETARGETENVIRONMENT
    // Open in the browser the employee chooses.
    NONE_BROWSERSITETARGETENVIRONMENT
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shouldn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSITETARGETENVIRONMENT
)

func (i BrowserSiteTargetEnvironment) String() string {
    return []string{"internetExplorerMode", "internetExplorer11", "microsoftEdge", "configurable", "none", "unknownFutureValue"}[i]
}
func ParseBrowserSiteTargetEnvironment(v string) (any, error) {
    result := INTERNETEXPLORERMODE_BROWSERSITETARGETENVIRONMENT
    switch v {
        case "internetExplorerMode":
            result = INTERNETEXPLORERMODE_BROWSERSITETARGETENVIRONMENT
        case "internetExplorer11":
            result = INTERNETEXPLORER11_BROWSERSITETARGETENVIRONMENT
        case "microsoftEdge":
            result = MICROSOFTEDGE_BROWSERSITETARGETENVIRONMENT
        case "configurable":
            result = CONFIGURABLE_BROWSERSITETARGETENVIRONMENT
        case "none":
            result = NONE_BROWSERSITETARGETENVIRONMENT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSITETARGETENVIRONMENT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSiteTargetEnvironment(values []BrowserSiteTargetEnvironment) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSiteTargetEnvironment) isMultiValue() bool {
    return false
}
