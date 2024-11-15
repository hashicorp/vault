package models
type SearchAlterationType int

const (
    SUGGESTION_SEARCHALTERATIONTYPE SearchAlterationType = iota
    MODIFICATION_SEARCHALTERATIONTYPE
    UNKNOWNFUTUREVALUE_SEARCHALTERATIONTYPE
)

func (i SearchAlterationType) String() string {
    return []string{"suggestion", "modification", "unknownFutureValue"}[i]
}
func ParseSearchAlterationType(v string) (any, error) {
    result := SUGGESTION_SEARCHALTERATIONTYPE
    switch v {
        case "suggestion":
            result = SUGGESTION_SEARCHALTERATIONTYPE
        case "modification":
            result = MODIFICATION_SEARCHALTERATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SEARCHALTERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSearchAlterationType(values []SearchAlterationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SearchAlterationType) isMultiValue() bool {
    return false
}
