package models
type VirtualEventAttendeeRegistrationStatus int

const (
    REGISTERED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS VirtualEventAttendeeRegistrationStatus = iota
    CANCELED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
    WAITLISTED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
    PENDINGAPPROVAL_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
    REJECTEDBYORGANIZER_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
    UNKNOWNFUTUREVALUE_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
)

func (i VirtualEventAttendeeRegistrationStatus) String() string {
    return []string{"registered", "canceled", "waitlisted", "pendingApproval", "rejectedByOrganizer", "unknownFutureValue"}[i]
}
func ParseVirtualEventAttendeeRegistrationStatus(v string) (any, error) {
    result := REGISTERED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
    switch v {
        case "registered":
            result = REGISTERED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        case "canceled":
            result = CANCELED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        case "waitlisted":
            result = WAITLISTED_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        case "pendingApproval":
            result = PENDINGAPPROVAL_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        case "rejectedByOrganizer":
            result = REJECTEDBYORGANIZER_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIRTUALEVENTATTENDEEREGISTRATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVirtualEventAttendeeRegistrationStatus(values []VirtualEventAttendeeRegistrationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VirtualEventAttendeeRegistrationStatus) isMultiValue() bool {
    return false
}
