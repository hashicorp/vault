package models
type VirtualAppointmentMessageType int

const (
    CONFIRMATION_VIRTUALAPPOINTMENTMESSAGETYPE VirtualAppointmentMessageType = iota
    RESCHEDULE_VIRTUALAPPOINTMENTMESSAGETYPE
    CANCELLATION_VIRTUALAPPOINTMENTMESSAGETYPE
    UNKNOWNFUTUREVALUE_VIRTUALAPPOINTMENTMESSAGETYPE
)

func (i VirtualAppointmentMessageType) String() string {
    return []string{"confirmation", "reschedule", "cancellation", "unknownFutureValue"}[i]
}
func ParseVirtualAppointmentMessageType(v string) (any, error) {
    result := CONFIRMATION_VIRTUALAPPOINTMENTMESSAGETYPE
    switch v {
        case "confirmation":
            result = CONFIRMATION_VIRTUALAPPOINTMENTMESSAGETYPE
        case "reschedule":
            result = RESCHEDULE_VIRTUALAPPOINTMENTMESSAGETYPE
        case "cancellation":
            result = CANCELLATION_VIRTUALAPPOINTMENTMESSAGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIRTUALAPPOINTMENTMESSAGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVirtualAppointmentMessageType(values []VirtualAppointmentMessageType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VirtualAppointmentMessageType) isMultiValue() bool {
    return false
}
