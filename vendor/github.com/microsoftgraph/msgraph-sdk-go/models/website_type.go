package models
type WebsiteType int

const (
    OTHER_WEBSITETYPE WebsiteType = iota
    HOME_WEBSITETYPE
    WORK_WEBSITETYPE
    BLOG_WEBSITETYPE
    PROFILE_WEBSITETYPE
)

func (i WebsiteType) String() string {
    return []string{"other", "home", "work", "blog", "profile"}[i]
}
func ParseWebsiteType(v string) (any, error) {
    result := OTHER_WEBSITETYPE
    switch v {
        case "other":
            result = OTHER_WEBSITETYPE
        case "home":
            result = HOME_WEBSITETYPE
        case "work":
            result = WORK_WEBSITETYPE
        case "blog":
            result = BLOG_WEBSITETYPE
        case "profile":
            result = PROFILE_WEBSITETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWebsiteType(values []WebsiteType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WebsiteType) isMultiValue() bool {
    return false
}
