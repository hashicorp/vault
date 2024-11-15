package models
type BookingsServiceAvailabilityType int

const (
    BOOKWHENSTAFFAREFREE_BOOKINGSSERVICEAVAILABILITYTYPE BookingsServiceAvailabilityType = iota
    NOTBOOKABLE_BOOKINGSSERVICEAVAILABILITYTYPE
    CUSTOMWEEKLYHOURS_BOOKINGSSERVICEAVAILABILITYTYPE
    UNKNOWNFUTUREVALUE_BOOKINGSSERVICEAVAILABILITYTYPE
)

func (i BookingsServiceAvailabilityType) String() string {
    return []string{"bookWhenStaffAreFree", "notBookable", "customWeeklyHours", "unknownFutureValue"}[i]
}
func ParseBookingsServiceAvailabilityType(v string) (any, error) {
    result := BOOKWHENSTAFFAREFREE_BOOKINGSSERVICEAVAILABILITYTYPE
    switch v {
        case "bookWhenStaffAreFree":
            result = BOOKWHENSTAFFAREFREE_BOOKINGSSERVICEAVAILABILITYTYPE
        case "notBookable":
            result = NOTBOOKABLE_BOOKINGSSERVICEAVAILABILITYTYPE
        case "customWeeklyHours":
            result = CUSTOMWEEKLYHOURS_BOOKINGSSERVICEAVAILABILITYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGSSERVICEAVAILABILITYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingsServiceAvailabilityType(values []BookingsServiceAvailabilityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingsServiceAvailabilityType) isMultiValue() bool {
    return false
}
