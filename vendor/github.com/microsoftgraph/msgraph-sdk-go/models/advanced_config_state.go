package models
type AdvancedConfigState int

const (
    DEFAULT_ADVANCEDCONFIGSTATE AdvancedConfigState = iota
    ENABLED_ADVANCEDCONFIGSTATE
    DISABLED_ADVANCEDCONFIGSTATE
    UNKNOWNFUTUREVALUE_ADVANCEDCONFIGSTATE
)

func (i AdvancedConfigState) String() string {
    return []string{"default", "enabled", "disabled", "unknownFutureValue"}[i]
}
func ParseAdvancedConfigState(v string) (any, error) {
    result := DEFAULT_ADVANCEDCONFIGSTATE
    switch v {
        case "default":
            result = DEFAULT_ADVANCEDCONFIGSTATE
        case "enabled":
            result = ENABLED_ADVANCEDCONFIGSTATE
        case "disabled":
            result = DISABLED_ADVANCEDCONFIGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ADVANCEDCONFIGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAdvancedConfigState(values []AdvancedConfigState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AdvancedConfigState) isMultiValue() bool {
    return false
}
