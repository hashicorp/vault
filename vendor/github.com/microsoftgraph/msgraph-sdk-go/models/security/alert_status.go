package security
type AlertStatus int

const (
    UNKNOWN_ALERTSTATUS AlertStatus = iota
    NEW_ALERTSTATUS
    INPROGRESS_ALERTSTATUS
    RESOLVED_ALERTSTATUS
    UNKNOWNFUTUREVALUE_ALERTSTATUS
)

func (i AlertStatus) String() string {
    return []string{"unknown", "new", "inProgress", "resolved", "unknownFutureValue"}[i]
}
func ParseAlertStatus(v string) (any, error) {
    result := UNKNOWN_ALERTSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_ALERTSTATUS
        case "new":
            result = NEW_ALERTSTATUS
        case "inProgress":
            result = INPROGRESS_ALERTSTATUS
        case "resolved":
            result = RESOLVED_ALERTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ALERTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAlertStatus(values []AlertStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AlertStatus) isMultiValue() bool {
    return false
}
