package models
type BrowserSiteListStatus int

const (
    // A site list that has not yet been published
    DRAFT_BROWSERSITELISTSTATUS BrowserSiteListStatus = iota
    // A site list that has been published with no pending changes.
    PUBLISHED_BROWSERSITELISTSTATUS
    // A site that has pending changes
    PENDING_BROWSERSITELISTSTATUS
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shoudn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSITELISTSTATUS
)

func (i BrowserSiteListStatus) String() string {
    return []string{"draft", "published", "pending", "unknownFutureValue"}[i]
}
func ParseBrowserSiteListStatus(v string) (any, error) {
    result := DRAFT_BROWSERSITELISTSTATUS
    switch v {
        case "draft":
            result = DRAFT_BROWSERSITELISTSTATUS
        case "published":
            result = PUBLISHED_BROWSERSITELISTSTATUS
        case "pending":
            result = PENDING_BROWSERSITELISTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSITELISTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSiteListStatus(values []BrowserSiteListStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSiteListStatus) isMultiValue() bool {
    return false
}
