package models
type AttributeMappingSourceType int

const (
    ATTRIBUTE_ATTRIBUTEMAPPINGSOURCETYPE AttributeMappingSourceType = iota
    CONSTANT_ATTRIBUTEMAPPINGSOURCETYPE
    FUNCTION_ATTRIBUTEMAPPINGSOURCETYPE
)

func (i AttributeMappingSourceType) String() string {
    return []string{"Attribute", "Constant", "Function"}[i]
}
func ParseAttributeMappingSourceType(v string) (any, error) {
    result := ATTRIBUTE_ATTRIBUTEMAPPINGSOURCETYPE
    switch v {
        case "Attribute":
            result = ATTRIBUTE_ATTRIBUTEMAPPINGSOURCETYPE
        case "Constant":
            result = CONSTANT_ATTRIBUTEMAPPINGSOURCETYPE
        case "Function":
            result = FUNCTION_ATTRIBUTEMAPPINGSOURCETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAttributeMappingSourceType(values []AttributeMappingSourceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AttributeMappingSourceType) isMultiValue() bool {
    return false
}
