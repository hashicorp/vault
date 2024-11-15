package models
type BookingType int

const (
    UNKNOWN_BOOKINGTYPE BookingType = iota
    STANDARD_BOOKINGTYPE
    RESERVED_BOOKINGTYPE
)

func (i BookingType) String() string {
    return []string{"unknown", "standard", "reserved"}[i]
}
func ParseBookingType(v string) (any, error) {
    result := UNKNOWN_BOOKINGTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_BOOKINGTYPE
        case "standard":
            result = STANDARD_BOOKINGTYPE
        case "reserved":
            result = RESERVED_BOOKINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingType(values []BookingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingType) isMultiValue() bool {
    return false
}
