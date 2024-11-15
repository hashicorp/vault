package models
// Windows Hello for Business pin usage options
type WindowsHelloForBusinessPinUsage int

const (
    // Allowed the usage of certain pin rule
    ALLOWED_WINDOWSHELLOFORBUSINESSPINUSAGE WindowsHelloForBusinessPinUsage = iota
    // Enforce the usage of certain pin rule
    REQUIRED_WINDOWSHELLOFORBUSINESSPINUSAGE
    // Forbit the usage of certain pin rule
    DISALLOWED_WINDOWSHELLOFORBUSINESSPINUSAGE
)

func (i WindowsHelloForBusinessPinUsage) String() string {
    return []string{"allowed", "required", "disallowed"}[i]
}
func ParseWindowsHelloForBusinessPinUsage(v string) (any, error) {
    result := ALLOWED_WINDOWSHELLOFORBUSINESSPINUSAGE
    switch v {
        case "allowed":
            result = ALLOWED_WINDOWSHELLOFORBUSINESSPINUSAGE
        case "required":
            result = REQUIRED_WINDOWSHELLOFORBUSINESSPINUSAGE
        case "disallowed":
            result = DISALLOWED_WINDOWSHELLOFORBUSINESSPINUSAGE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsHelloForBusinessPinUsage(values []WindowsHelloForBusinessPinUsage) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsHelloForBusinessPinUsage) isMultiValue() bool {
    return false
}
