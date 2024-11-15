package security
type IncidentStatus int

const (
    ACTIVE_INCIDENTSTATUS IncidentStatus = iota
    RESOLVED_INCIDENTSTATUS
    INPROGRESS_INCIDENTSTATUS
    REDIRECTED_INCIDENTSTATUS
    UNKNOWNFUTUREVALUE_INCIDENTSTATUS
    AWAITINGACTION_INCIDENTSTATUS
)

func (i IncidentStatus) String() string {
    return []string{"active", "resolved", "inProgress", "redirected", "unknownFutureValue", "awaitingAction"}[i]
}
func ParseIncidentStatus(v string) (any, error) {
    result := ACTIVE_INCIDENTSTATUS
    switch v {
        case "active":
            result = ACTIVE_INCIDENTSTATUS
        case "resolved":
            result = RESOLVED_INCIDENTSTATUS
        case "inProgress":
            result = INPROGRESS_INCIDENTSTATUS
        case "redirected":
            result = REDIRECTED_INCIDENTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INCIDENTSTATUS
        case "awaitingAction":
            result = AWAITINGACTION_INCIDENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIncidentStatus(values []IncidentStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IncidentStatus) isMultiValue() bool {
    return false
}
