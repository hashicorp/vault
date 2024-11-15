package models
type BrowserSiteMergeType int

const (
    // No merge type
    NOMERGE_BROWSERSITEMERGETYPE BrowserSiteMergeType = iota
    // Default merge type
    DEFAULT_BROWSERSITEMERGETYPE
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shouldn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSITEMERGETYPE
)

func (i BrowserSiteMergeType) String() string {
    return []string{"noMerge", "default", "unknownFutureValue"}[i]
}
func ParseBrowserSiteMergeType(v string) (any, error) {
    result := NOMERGE_BROWSERSITEMERGETYPE
    switch v {
        case "noMerge":
            result = NOMERGE_BROWSERSITEMERGETYPE
        case "default":
            result = DEFAULT_BROWSERSITEMERGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSITEMERGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSiteMergeType(values []BrowserSiteMergeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSiteMergeType) isMultiValue() bool {
    return false
}
