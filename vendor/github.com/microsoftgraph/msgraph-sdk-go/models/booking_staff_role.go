package models
type BookingStaffRole int

const (
    GUEST_BOOKINGSTAFFROLE BookingStaffRole = iota
    ADMINISTRATOR_BOOKINGSTAFFROLE
    VIEWER_BOOKINGSTAFFROLE
    EXTERNALGUEST_BOOKINGSTAFFROLE
    UNKNOWNFUTUREVALUE_BOOKINGSTAFFROLE
    SCHEDULER_BOOKINGSTAFFROLE
    TEAMMEMBER_BOOKINGSTAFFROLE
)

func (i BookingStaffRole) String() string {
    return []string{"guest", "administrator", "viewer", "externalGuest", "unknownFutureValue", "scheduler", "teamMember"}[i]
}
func ParseBookingStaffRole(v string) (any, error) {
    result := GUEST_BOOKINGSTAFFROLE
    switch v {
        case "guest":
            result = GUEST_BOOKINGSTAFFROLE
        case "administrator":
            result = ADMINISTRATOR_BOOKINGSTAFFROLE
        case "viewer":
            result = VIEWER_BOOKINGSTAFFROLE
        case "externalGuest":
            result = EXTERNALGUEST_BOOKINGSTAFFROLE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGSTAFFROLE
        case "scheduler":
            result = SCHEDULER_BOOKINGSTAFFROLE
        case "teamMember":
            result = TEAMMEMBER_BOOKINGSTAFFROLE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingStaffRole(values []BookingStaffRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingStaffRole) isMultiValue() bool {
    return false
}
