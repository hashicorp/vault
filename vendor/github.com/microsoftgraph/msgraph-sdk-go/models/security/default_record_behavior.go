package security
type DefaultRecordBehavior int

const (
    STARTLOCKED_DEFAULTRECORDBEHAVIOR DefaultRecordBehavior = iota
    STARTUNLOCKED_DEFAULTRECORDBEHAVIOR
    UNKNOWNFUTUREVALUE_DEFAULTRECORDBEHAVIOR
)

func (i DefaultRecordBehavior) String() string {
    return []string{"startLocked", "startUnlocked", "unknownFutureValue"}[i]
}
func ParseDefaultRecordBehavior(v string) (any, error) {
    result := STARTLOCKED_DEFAULTRECORDBEHAVIOR
    switch v {
        case "startLocked":
            result = STARTLOCKED_DEFAULTRECORDBEHAVIOR
        case "startUnlocked":
            result = STARTUNLOCKED_DEFAULTRECORDBEHAVIOR
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEFAULTRECORDBEHAVIOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDefaultRecordBehavior(values []DefaultRecordBehavior) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DefaultRecordBehavior) isMultiValue() bool {
    return false
}
