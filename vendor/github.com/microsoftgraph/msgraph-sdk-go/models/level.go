package models
type Level int

const (
    BEGINNER_LEVEL Level = iota
    INTERMEDIATE_LEVEL
    ADVANCED_LEVEL
    UNKNOWNFUTUREVALUE_LEVEL
)

func (i Level) String() string {
    return []string{"beginner", "intermediate", "advanced", "unknownFutureValue"}[i]
}
func ParseLevel(v string) (any, error) {
    result := BEGINNER_LEVEL
    switch v {
        case "beginner":
            result = BEGINNER_LEVEL
        case "intermediate":
            result = INTERMEDIATE_LEVEL
        case "advanced":
            result = ADVANCED_LEVEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLevel(values []Level) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Level) isMultiValue() bool {
    return false
}
