package models
// Character set which is to be used for a user's app PIN
type ManagedAppPinCharacterSet int

const (
    // Numeric characters
    NUMERIC_MANAGEDAPPPINCHARACTERSET ManagedAppPinCharacterSet = iota
    // Alphanumeric and symbolic characters
    ALPHANUMERICANDSYMBOL_MANAGEDAPPPINCHARACTERSET
)

func (i ManagedAppPinCharacterSet) String() string {
    return []string{"numeric", "alphanumericAndSymbol"}[i]
}
func ParseManagedAppPinCharacterSet(v string) (any, error) {
    result := NUMERIC_MANAGEDAPPPINCHARACTERSET
    switch v {
        case "numeric":
            result = NUMERIC_MANAGEDAPPPINCHARACTERSET
        case "alphanumericAndSymbol":
            result = ALPHANUMERICANDSYMBOL_MANAGEDAPPPINCHARACTERSET
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedAppPinCharacterSet(values []ManagedAppPinCharacterSet) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedAppPinCharacterSet) isMultiValue() bool {
    return false
}
