package models
type PrintEvent int

const (
    JOBSTARTED_PRINTEVENT PrintEvent = iota
    UNKNOWNFUTUREVALUE_PRINTEVENT
)

func (i PrintEvent) String() string {
    return []string{"jobStarted", "unknownFutureValue"}[i]
}
func ParsePrintEvent(v string) (any, error) {
    result := JOBSTARTED_PRINTEVENT
    switch v {
        case "jobStarted":
            result = JOBSTARTED_PRINTEVENT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTEVENT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintEvent(values []PrintEvent) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintEvent) isMultiValue() bool {
    return false
}
