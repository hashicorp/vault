package models
type FilterMode int

const (
    INCLUDE_FILTERMODE FilterMode = iota
    EXCLUDE_FILTERMODE
)

func (i FilterMode) String() string {
    return []string{"include", "exclude"}[i]
}
func ParseFilterMode(v string) (any, error) {
    result := INCLUDE_FILTERMODE
    switch v {
        case "include":
            result = INCLUDE_FILTERMODE
        case "exclude":
            result = EXCLUDE_FILTERMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFilterMode(values []FilterMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FilterMode) isMultiValue() bool {
    return false
}
