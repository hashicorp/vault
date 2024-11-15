package models
type ConditionalAccessClientApp int

const (
    ALL_CONDITIONALACCESSCLIENTAPP ConditionalAccessClientApp = iota
    BROWSER_CONDITIONALACCESSCLIENTAPP
    MOBILEAPPSANDDESKTOPCLIENTS_CONDITIONALACCESSCLIENTAPP
    EXCHANGEACTIVESYNC_CONDITIONALACCESSCLIENTAPP
    EASSUPPORTED_CONDITIONALACCESSCLIENTAPP
    OTHER_CONDITIONALACCESSCLIENTAPP
    UNKNOWNFUTUREVALUE_CONDITIONALACCESSCLIENTAPP
)

func (i ConditionalAccessClientApp) String() string {
    return []string{"all", "browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "easSupported", "other", "unknownFutureValue"}[i]
}
func ParseConditionalAccessClientApp(v string) (any, error) {
    result := ALL_CONDITIONALACCESSCLIENTAPP
    switch v {
        case "all":
            result = ALL_CONDITIONALACCESSCLIENTAPP
        case "browser":
            result = BROWSER_CONDITIONALACCESSCLIENTAPP
        case "mobileAppsAndDesktopClients":
            result = MOBILEAPPSANDDESKTOPCLIENTS_CONDITIONALACCESSCLIENTAPP
        case "exchangeActiveSync":
            result = EXCHANGEACTIVESYNC_CONDITIONALACCESSCLIENTAPP
        case "easSupported":
            result = EASSUPPORTED_CONDITIONALACCESSCLIENTAPP
        case "other":
            result = OTHER_CONDITIONALACCESSCLIENTAPP
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONDITIONALACCESSCLIENTAPP
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeConditionalAccessClientApp(values []ConditionalAccessClientApp) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConditionalAccessClientApp) isMultiValue() bool {
    return false
}
