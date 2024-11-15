package models
type BrowserSiteStatus int

const (
    // A site that has been published
    PUBLISHED_BROWSERSITESTATUS BrowserSiteStatus = iota
    // A site that has been added pending publish
    PENDINGADD_BROWSERSITESTATUS
    // A site that has been edited pending publish
    PENDINGEDIT_BROWSERSITESTATUS
    // A site that has been deleted pending publish
    PENDINGDELETE_BROWSERSITESTATUS
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shouldn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSITESTATUS
)

func (i BrowserSiteStatus) String() string {
    return []string{"published", "pendingAdd", "pendingEdit", "pendingDelete", "unknownFutureValue"}[i]
}
func ParseBrowserSiteStatus(v string) (any, error) {
    result := PUBLISHED_BROWSERSITESTATUS
    switch v {
        case "published":
            result = PUBLISHED_BROWSERSITESTATUS
        case "pendingAdd":
            result = PENDINGADD_BROWSERSITESTATUS
        case "pendingEdit":
            result = PENDINGEDIT_BROWSERSITESTATUS
        case "pendingDelete":
            result = PENDINGDELETE_BROWSERSITESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSITESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSiteStatus(values []BrowserSiteStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSiteStatus) isMultiValue() bool {
    return false
}
