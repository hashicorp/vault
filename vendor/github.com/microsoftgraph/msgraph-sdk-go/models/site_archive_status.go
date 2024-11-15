package models
type SiteArchiveStatus int

const (
    RECENTLYARCHIVED_SITEARCHIVESTATUS SiteArchiveStatus = iota
    FULLYARCHIVED_SITEARCHIVESTATUS
    REACTIVATING_SITEARCHIVESTATUS
    UNKNOWNFUTUREVALUE_SITEARCHIVESTATUS
)

func (i SiteArchiveStatus) String() string {
    return []string{"recentlyArchived", "fullyArchived", "reactivating", "unknownFutureValue"}[i]
}
func ParseSiteArchiveStatus(v string) (any, error) {
    result := RECENTLYARCHIVED_SITEARCHIVESTATUS
    switch v {
        case "recentlyArchived":
            result = RECENTLYARCHIVED_SITEARCHIVESTATUS
        case "fullyArchived":
            result = FULLYARCHIVED_SITEARCHIVESTATUS
        case "reactivating":
            result = REACTIVATING_SITEARCHIVESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SITEARCHIVESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSiteArchiveStatus(values []SiteArchiveStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SiteArchiveStatus) isMultiValue() bool {
    return false
}
