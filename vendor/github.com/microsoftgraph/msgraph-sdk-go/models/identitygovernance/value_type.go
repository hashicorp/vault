package identitygovernance
type ValueType int

const (
    ENUM_VALUETYPE ValueType = iota
    STRING_VALUETYPE
    INT_VALUETYPE
    BOOL_VALUETYPE
    UNKNOWNFUTUREVALUE_VALUETYPE
)

func (i ValueType) String() string {
    return []string{"enum", "string", "int", "bool", "unknownFutureValue"}[i]
}
func ParseValueType(v string) (any, error) {
    result := ENUM_VALUETYPE
    switch v {
        case "enum":
            result = ENUM_VALUETYPE
        case "string":
            result = STRING_VALUETYPE
        case "int":
            result = INT_VALUETYPE
        case "bool":
            result = BOOL_VALUETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VALUETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeValueType(values []ValueType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ValueType) isMultiValue() bool {
    return false
}
