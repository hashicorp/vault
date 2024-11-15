package externalconnectors
type AccessType int

const (
    GRANT_ACCESSTYPE AccessType = iota
    DENY_ACCESSTYPE
    UNKNOWNFUTUREVALUE_ACCESSTYPE
)

func (i AccessType) String() string {
    return []string{"grant", "deny", "unknownFutureValue"}[i]
}
func ParseAccessType(v string) (any, error) {
    result := GRANT_ACCESSTYPE
    switch v {
        case "grant":
            result = GRANT_ACCESSTYPE
        case "deny":
            result = DENY_ACCESSTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCESSTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccessType(values []AccessType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccessType) isMultiValue() bool {
    return false
}
