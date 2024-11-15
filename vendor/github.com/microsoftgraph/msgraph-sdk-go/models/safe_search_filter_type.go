package models
// Specifies what level of safe search (filtering adult content) is required
type SafeSearchFilterType int

const (
    // User Defined, default value, no intent.
    USERDEFINED_SAFESEARCHFILTERTYPE SafeSearchFilterType = iota
    // Strict, highest filtering against adult content.
    STRICT_SAFESEARCHFILTERTYPE
    // Moderate filtering against adult content (valid search results will not be filtered).
    MODERATE_SAFESEARCHFILTERTYPE
)

func (i SafeSearchFilterType) String() string {
    return []string{"userDefined", "strict", "moderate"}[i]
}
func ParseSafeSearchFilterType(v string) (any, error) {
    result := USERDEFINED_SAFESEARCHFILTERTYPE
    switch v {
        case "userDefined":
            result = USERDEFINED_SAFESEARCHFILTERTYPE
        case "strict":
            result = STRICT_SAFESEARCHFILTERTYPE
        case "moderate":
            result = MODERATE_SAFESEARCHFILTERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSafeSearchFilterType(values []SafeSearchFilterType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SafeSearchFilterType) isMultiValue() bool {
    return false
}
