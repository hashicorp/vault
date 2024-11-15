package models
type MediaState int

const (
    ACTIVE_MEDIASTATE MediaState = iota
    INACTIVE_MEDIASTATE
    UNKNOWNFUTUREVALUE_MEDIASTATE
)

func (i MediaState) String() string {
    return []string{"active", "inactive", "unknownFutureValue"}[i]
}
func ParseMediaState(v string) (any, error) {
    result := ACTIVE_MEDIASTATE
    switch v {
        case "active":
            result = ACTIVE_MEDIASTATE
        case "inactive":
            result = INACTIVE_MEDIASTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MEDIASTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMediaState(values []MediaState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MediaState) isMultiValue() bool {
    return false
}
