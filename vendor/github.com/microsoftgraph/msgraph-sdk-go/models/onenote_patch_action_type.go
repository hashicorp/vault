package models
type OnenotePatchActionType int

const (
    REPLACE_ONENOTEPATCHACTIONTYPE OnenotePatchActionType = iota
    APPEND_ONENOTEPATCHACTIONTYPE
    DELETE_ONENOTEPATCHACTIONTYPE
    INSERT_ONENOTEPATCHACTIONTYPE
    PREPEND_ONENOTEPATCHACTIONTYPE
)

func (i OnenotePatchActionType) String() string {
    return []string{"Replace", "Append", "Delete", "Insert", "Prepend"}[i]
}
func ParseOnenotePatchActionType(v string) (any, error) {
    result := REPLACE_ONENOTEPATCHACTIONTYPE
    switch v {
        case "Replace":
            result = REPLACE_ONENOTEPATCHACTIONTYPE
        case "Append":
            result = APPEND_ONENOTEPATCHACTIONTYPE
        case "Delete":
            result = DELETE_ONENOTEPATCHACTIONTYPE
        case "Insert":
            result = INSERT_ONENOTEPATCHACTIONTYPE
        case "Prepend":
            result = PREPEND_ONENOTEPATCHACTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOnenotePatchActionType(values []OnenotePatchActionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnenotePatchActionType) isMultiValue() bool {
    return false
}
