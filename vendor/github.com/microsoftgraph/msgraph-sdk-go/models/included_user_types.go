package models
type IncludedUserTypes int

const (
    ALL_INCLUDEDUSERTYPES IncludedUserTypes = iota
    MEMBER_INCLUDEDUSERTYPES
    GUEST_INCLUDEDUSERTYPES
    UNKNOWNFUTUREVALUE_INCLUDEDUSERTYPES
)

func (i IncludedUserTypes) String() string {
    return []string{"all", "member", "guest", "unknownFutureValue"}[i]
}
func ParseIncludedUserTypes(v string) (any, error) {
    result := ALL_INCLUDEDUSERTYPES
    switch v {
        case "all":
            result = ALL_INCLUDEDUSERTYPES
        case "member":
            result = MEMBER_INCLUDEDUSERTYPES
        case "guest":
            result = GUEST_INCLUDEDUSERTYPES
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INCLUDEDUSERTYPES
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIncludedUserTypes(values []IncludedUserTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IncludedUserTypes) isMultiValue() bool {
    return false
}
