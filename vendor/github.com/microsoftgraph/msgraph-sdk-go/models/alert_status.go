package models
type AlertStatus int

const (
    UNKNOWN_ALERTSTATUS AlertStatus = iota
    NEWALERT_ALERTSTATUS
    INPROGRESS_ALERTSTATUS
    RESOLVED_ALERTSTATUS
    DISMISSED_ALERTSTATUS
    UNKNOWNFUTUREVALUE_ALERTSTATUS
)

func (i AlertStatus) String() string {
    return []string{"unknown", "newAlert", "inProgress", "resolved", "dismissed", "unknownFutureValue"}[i]
}
func ParseAlertStatus(v string) (any, error) {
    result := UNKNOWN_ALERTSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_ALERTSTATUS
        case "newAlert":
            result = NEWALERT_ALERTSTATUS
        case "inProgress":
            result = INPROGRESS_ALERTSTATUS
        case "resolved":
            result = RESOLVED_ALERTSTATUS
        case "dismissed":
            result = DISMISSED_ALERTSTATUS
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
