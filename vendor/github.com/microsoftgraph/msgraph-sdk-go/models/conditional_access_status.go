package models
type ConditionalAccessStatus int

const (
    SUCCESS_CONDITIONALACCESSSTATUS ConditionalAccessStatus = iota
    FAILURE_CONDITIONALACCESSSTATUS
    NOTAPPLIED_CONDITIONALACCESSSTATUS
    UNKNOWNFUTUREVALUE_CONDITIONALACCESSSTATUS
)

func (i ConditionalAccessStatus) String() string {
    return []string{"success", "failure", "notApplied", "unknownFutureValue"}[i]
}
func ParseConditionalAccessStatus(v string) (any, error) {
    result := SUCCESS_CONDITIONALACCESSSTATUS
    switch v {
        case "success":
            result = SUCCESS_CONDITIONALACCESSSTATUS
        case "failure":
            result = FAILURE_CONDITIONALACCESSSTATUS
        case "notApplied":
            result = NOTAPPLIED_CONDITIONALACCESSSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONDITIONALACCESSSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeConditionalAccessStatus(values []ConditionalAccessStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConditionalAccessStatus) isMultiValue() bool {
    return false
}
