package models
// Contains properties for detection operator.
type Win32LobAppRuleOperator int

const (
    // Not configured.
    NOTCONFIGURED_WIN32LOBAPPRULEOPERATOR Win32LobAppRuleOperator = iota
    // Equal operator.
    EQUAL_WIN32LOBAPPRULEOPERATOR
    // Not equal operator.
    NOTEQUAL_WIN32LOBAPPRULEOPERATOR
    // Greater than operator.
    GREATERTHAN_WIN32LOBAPPRULEOPERATOR
    // Greater than or equal operator.
    GREATERTHANOREQUAL_WIN32LOBAPPRULEOPERATOR
    // Less than operator.
    LESSTHAN_WIN32LOBAPPRULEOPERATOR
    // Less than or equal operator.
    LESSTHANOREQUAL_WIN32LOBAPPRULEOPERATOR
)

func (i Win32LobAppRuleOperator) String() string {
    return []string{"notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"}[i]
}
func ParseWin32LobAppRuleOperator(v string) (any, error) {
    result := NOTCONFIGURED_WIN32LOBAPPRULEOPERATOR
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WIN32LOBAPPRULEOPERATOR
        case "equal":
            result = EQUAL_WIN32LOBAPPRULEOPERATOR
        case "notEqual":
            result = NOTEQUAL_WIN32LOBAPPRULEOPERATOR
        case "greaterThan":
            result = GREATERTHAN_WIN32LOBAPPRULEOPERATOR
        case "greaterThanOrEqual":
            result = GREATERTHANOREQUAL_WIN32LOBAPPRULEOPERATOR
        case "lessThan":
            result = LESSTHAN_WIN32LOBAPPRULEOPERATOR
        case "lessThanOrEqual":
            result = LESSTHANOREQUAL_WIN32LOBAPPRULEOPERATOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppRuleOperator(values []Win32LobAppRuleOperator) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppRuleOperator) isMultiValue() bool {
    return false
}
