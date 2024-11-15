package models
type Status int

const (
    ACTIVE_STATUS Status = iota
    UPDATED_STATUS
    DELETED_STATUS
    IGNORED_STATUS
    UNKNOWNFUTUREVALUE_STATUS
)

func (i Status) String() string {
    return []string{"active", "updated", "deleted", "ignored", "unknownFutureValue"}[i]
}
func ParseStatus(v string) (any, error) {
    result := ACTIVE_STATUS
    switch v {
        case "active":
            result = ACTIVE_STATUS
        case "updated":
            result = UPDATED_STATUS
        case "deleted":
            result = DELETED_STATUS
        case "ignored":
            result = IGNORED_STATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_STATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeStatus(values []Status) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Status) isMultiValue() bool {
    return false
}
