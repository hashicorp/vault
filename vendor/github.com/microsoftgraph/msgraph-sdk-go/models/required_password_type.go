package models
// Possible values of required passwords.
type RequiredPasswordType int

const (
    // Device default value, no intent.
    DEVICEDEFAULT_REQUIREDPASSWORDTYPE RequiredPasswordType = iota
    // Alphanumeric password required.
    ALPHANUMERIC_REQUIREDPASSWORDTYPE
    // Numeric password required.
    NUMERIC_REQUIREDPASSWORDTYPE
)

func (i RequiredPasswordType) String() string {
    return []string{"deviceDefault", "alphanumeric", "numeric"}[i]
}
func ParseRequiredPasswordType(v string) (any, error) {
    result := DEVICEDEFAULT_REQUIREDPASSWORDTYPE
    switch v {
        case "deviceDefault":
            result = DEVICEDEFAULT_REQUIREDPASSWORDTYPE
        case "alphanumeric":
            result = ALPHANUMERIC_REQUIREDPASSWORDTYPE
        case "numeric":
            result = NUMERIC_REQUIREDPASSWORDTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRequiredPasswordType(values []RequiredPasswordType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RequiredPasswordType) isMultiValue() bool {
    return false
}
