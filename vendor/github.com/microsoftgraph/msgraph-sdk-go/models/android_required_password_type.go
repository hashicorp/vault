package models
// Android required password type.
type AndroidRequiredPasswordType int

const (
    // Device default value, no intent.
    DEVICEDEFAULT_ANDROIDREQUIREDPASSWORDTYPE AndroidRequiredPasswordType = iota
    // Alphabetic password required.
    ALPHABETIC_ANDROIDREQUIREDPASSWORDTYPE
    // Alphanumeric password required.
    ALPHANUMERIC_ANDROIDREQUIREDPASSWORDTYPE
    // Alphanumeric with symbols password required.
    ALPHANUMERICWITHSYMBOLS_ANDROIDREQUIREDPASSWORDTYPE
    // Low security biometrics based password required.
    LOWSECURITYBIOMETRIC_ANDROIDREQUIREDPASSWORDTYPE
    // Numeric password required.
    NUMERIC_ANDROIDREQUIREDPASSWORDTYPE
    // Numeric complex password required.
    NUMERICCOMPLEX_ANDROIDREQUIREDPASSWORDTYPE
    // A password or pattern is required, and any is acceptable.
    ANY_ANDROIDREQUIREDPASSWORDTYPE
)

func (i AndroidRequiredPasswordType) String() string {
    return []string{"deviceDefault", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "numeric", "numericComplex", "any"}[i]
}
func ParseAndroidRequiredPasswordType(v string) (any, error) {
    result := DEVICEDEFAULT_ANDROIDREQUIREDPASSWORDTYPE
    switch v {
        case "deviceDefault":
            result = DEVICEDEFAULT_ANDROIDREQUIREDPASSWORDTYPE
        case "alphabetic":
            result = ALPHABETIC_ANDROIDREQUIREDPASSWORDTYPE
        case "alphanumeric":
            result = ALPHANUMERIC_ANDROIDREQUIREDPASSWORDTYPE
        case "alphanumericWithSymbols":
            result = ALPHANUMERICWITHSYMBOLS_ANDROIDREQUIREDPASSWORDTYPE
        case "lowSecurityBiometric":
            result = LOWSECURITYBIOMETRIC_ANDROIDREQUIREDPASSWORDTYPE
        case "numeric":
            result = NUMERIC_ANDROIDREQUIREDPASSWORDTYPE
        case "numericComplex":
            result = NUMERICCOMPLEX_ANDROIDREQUIREDPASSWORDTYPE
        case "any":
            result = ANY_ANDROIDREQUIREDPASSWORDTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAndroidRequiredPasswordType(values []AndroidRequiredPasswordType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AndroidRequiredPasswordType) isMultiValue() bool {
    return false
}
