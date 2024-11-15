package models
type RemindBeforeTimeInMinutesType int

const (
    MINS15_REMINDBEFORETIMEINMINUTESTYPE RemindBeforeTimeInMinutesType = iota
    UNKNOWNFUTUREVALUE_REMINDBEFORETIMEINMINUTESTYPE
)

func (i RemindBeforeTimeInMinutesType) String() string {
    return []string{"mins15", "unknownFutureValue"}[i]
}
func ParseRemindBeforeTimeInMinutesType(v string) (any, error) {
    result := MINS15_REMINDBEFORETIMEINMINUTESTYPE
    switch v {
        case "mins15":
            result = MINS15_REMINDBEFORETIMEINMINUTESTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_REMINDBEFORETIMEINMINUTESTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRemindBeforeTimeInMinutesType(values []RemindBeforeTimeInMinutesType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RemindBeforeTimeInMinutesType) isMultiValue() bool {
    return false
}
