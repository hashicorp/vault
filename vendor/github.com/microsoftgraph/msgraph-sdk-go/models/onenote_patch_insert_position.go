package models
type OnenotePatchInsertPosition int

const (
    AFTER_ONENOTEPATCHINSERTPOSITION OnenotePatchInsertPosition = iota
    BEFORE_ONENOTEPATCHINSERTPOSITION
)

func (i OnenotePatchInsertPosition) String() string {
    return []string{"After", "Before"}[i]
}
func ParseOnenotePatchInsertPosition(v string) (any, error) {
    result := AFTER_ONENOTEPATCHINSERTPOSITION
    switch v {
        case "After":
            result = AFTER_ONENOTEPATCHINSERTPOSITION
        case "Before":
            result = BEFORE_ONENOTEPATCHINSERTPOSITION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOnenotePatchInsertPosition(values []OnenotePatchInsertPosition) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnenotePatchInsertPosition) isMultiValue() bool {
    return false
}
