package models
type BrowserSharedCookieSourceEnvironment int

const (
    // Share session cookies from Microsoft Edge to Internet Explorer.
    MICROSOFTEDGE_BROWSERSHAREDCOOKIESOURCEENVIRONMENT BrowserSharedCookieSourceEnvironment = iota
    // Share session cookies from Internet Explorer to Microsoft Edge.
    INTERNETEXPLORER11_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
    // Share session cookies to and from Microsoft Edge and Internet Explorer.
    BOTH_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
    // Placeholder for evolvable enum, but this enum is never returned to the caller, so it shouldn't be necessary.
    UNKNOWNFUTUREVALUE_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
)

func (i BrowserSharedCookieSourceEnvironment) String() string {
    return []string{"microsoftEdge", "internetExplorer11", "both", "unknownFutureValue"}[i]
}
func ParseBrowserSharedCookieSourceEnvironment(v string) (any, error) {
    result := MICROSOFTEDGE_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
    switch v {
        case "microsoftEdge":
            result = MICROSOFTEDGE_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
        case "internetExplorer11":
            result = INTERNETEXPLORER11_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
        case "both":
            result = BOTH_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BROWSERSHAREDCOOKIESOURCEENVIRONMENT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBrowserSharedCookieSourceEnvironment(values []BrowserSharedCookieSourceEnvironment) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BrowserSharedCookieSourceEnvironment) isMultiValue() bool {
    return false
}
