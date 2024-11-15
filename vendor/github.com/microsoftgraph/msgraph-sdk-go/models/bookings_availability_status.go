package models
type BookingsAvailabilityStatus int

const (
    AVAILABLE_BOOKINGSAVAILABILITYSTATUS BookingsAvailabilityStatus = iota
    BUSY_BOOKINGSAVAILABILITYSTATUS
    SLOTSAVAILABLE_BOOKINGSAVAILABILITYSTATUS
    OUTOFOFFICE_BOOKINGSAVAILABILITYSTATUS
    UNKNOWNFUTUREVALUE_BOOKINGSAVAILABILITYSTATUS
)

func (i BookingsAvailabilityStatus) String() string {
    return []string{"available", "busy", "slotsAvailable", "outOfOffice", "unknownFutureValue"}[i]
}
func ParseBookingsAvailabilityStatus(v string) (any, error) {
    result := AVAILABLE_BOOKINGSAVAILABILITYSTATUS
    switch v {
        case "available":
            result = AVAILABLE_BOOKINGSAVAILABILITYSTATUS
        case "busy":
            result = BUSY_BOOKINGSAVAILABILITYSTATUS
        case "slotsAvailable":
            result = SLOTSAVAILABLE_BOOKINGSAVAILABILITYSTATUS
        case "outOfOffice":
            result = OUTOFOFFICE_BOOKINGSAVAILABILITYSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGSAVAILABILITYSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingsAvailabilityStatus(values []BookingsAvailabilityStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingsAvailabilityStatus) isMultiValue() bool {
    return false
}
