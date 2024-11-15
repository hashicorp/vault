package models
// Contains rule types for Win32 LOB apps.
type Win32LobAppRuleType int

const (
    // Detection rule.
    DETECTION_WIN32LOBAPPRULETYPE Win32LobAppRuleType = iota
    // Requirement rule.
    REQUIREMENT_WIN32LOBAPPRULETYPE
)

func (i Win32LobAppRuleType) String() string {
    return []string{"detection", "requirement"}[i]
}
func ParseWin32LobAppRuleType(v string) (any, error) {
    result := DETECTION_WIN32LOBAPPRULETYPE
    switch v {
        case "detection":
            result = DETECTION_WIN32LOBAPPRULETYPE
        case "requirement":
            result = REQUIREMENT_WIN32LOBAPPRULETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppRuleType(values []Win32LobAppRuleType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppRuleType) isMultiValue() bool {
    return false
}
