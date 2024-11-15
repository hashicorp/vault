package models
type WindowsUserType int

const (
    // Indicates that the user has administrator privileges.
    ADMINISTRATOR_WINDOWSUSERTYPE WindowsUserType = iota
    // Indicates that the user is a low-rights user without administrator privileges.
    STANDARD_WINDOWSUSERTYPE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_WINDOWSUSERTYPE
)

func (i WindowsUserType) String() string {
    return []string{"administrator", "standard", "unknownFutureValue"}[i]
}
func ParseWindowsUserType(v string) (any, error) {
    result := ADMINISTRATOR_WINDOWSUSERTYPE
    switch v {
        case "administrator":
            result = ADMINISTRATOR_WINDOWSUSERTYPE
        case "standard":
            result = STANDARD_WINDOWSUSERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WINDOWSUSERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsUserType(values []WindowsUserType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsUserType) isMultiValue() bool {
    return false
}
