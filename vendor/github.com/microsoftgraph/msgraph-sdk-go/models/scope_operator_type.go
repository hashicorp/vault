package models
type ScopeOperatorType int

const (
    BINARY_SCOPEOPERATORTYPE ScopeOperatorType = iota
    UNARY_SCOPEOPERATORTYPE
)

func (i ScopeOperatorType) String() string {
    return []string{"Binary", "Unary"}[i]
}
func ParseScopeOperatorType(v string) (any, error) {
    result := BINARY_SCOPEOPERATORTYPE
    switch v {
        case "Binary":
            result = BINARY_SCOPEOPERATORTYPE
        case "Unary":
            result = UNARY_SCOPEOPERATORTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeScopeOperatorType(values []ScopeOperatorType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ScopeOperatorType) isMultiValue() bool {
    return false
}
