package models
type MultiFactorAuthConfiguration int

const (
    NOTREQUIRED_MULTIFACTORAUTHCONFIGURATION MultiFactorAuthConfiguration = iota
    REQUIRED_MULTIFACTORAUTHCONFIGURATION
    UNKNOWNFUTUREVALUE_MULTIFACTORAUTHCONFIGURATION
)

func (i MultiFactorAuthConfiguration) String() string {
    return []string{"notRequired", "required", "unknownFutureValue"}[i]
}
func ParseMultiFactorAuthConfiguration(v string) (any, error) {
    result := NOTREQUIRED_MULTIFACTORAUTHCONFIGURATION
    switch v {
        case "notRequired":
            result = NOTREQUIRED_MULTIFACTORAUTHCONFIGURATION
        case "required":
            result = REQUIRED_MULTIFACTORAUTHCONFIGURATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MULTIFACTORAUTHCONFIGURATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMultiFactorAuthConfiguration(values []MultiFactorAuthConfiguration) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MultiFactorAuthConfiguration) isMultiValue() bool {
    return false
}
