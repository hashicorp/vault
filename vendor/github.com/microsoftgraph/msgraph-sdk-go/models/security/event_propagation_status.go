package security
type EventPropagationStatus int

const (
    NONE_EVENTPROPAGATIONSTATUS EventPropagationStatus = iota
    INPROCESSING_EVENTPROPAGATIONSTATUS
    FAILED_EVENTPROPAGATIONSTATUS
    SUCCESS_EVENTPROPAGATIONSTATUS
    UNKNOWNFUTUREVALUE_EVENTPROPAGATIONSTATUS
)

func (i EventPropagationStatus) String() string {
    return []string{"none", "inProcessing", "failed", "success", "unknownFutureValue"}[i]
}
func ParseEventPropagationStatus(v string) (any, error) {
    result := NONE_EVENTPROPAGATIONSTATUS
    switch v {
        case "none":
            result = NONE_EVENTPROPAGATIONSTATUS
        case "inProcessing":
            result = INPROCESSING_EVENTPROPAGATIONSTATUS
        case "failed":
            result = FAILED_EVENTPROPAGATIONSTATUS
        case "success":
            result = SUCCESS_EVENTPROPAGATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EVENTPROPAGATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEventPropagationStatus(values []EventPropagationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EventPropagationStatus) isMultiValue() bool {
    return false
}
