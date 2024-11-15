package models
type BookingPageAccessControl int

const (
    UNRESTRICTED_BOOKINGPAGEACCESSCONTROL BookingPageAccessControl = iota
    RESTRICTEDTOORGANIZATION_BOOKINGPAGEACCESSCONTROL
    UNKNOWNFUTUREVALUE_BOOKINGPAGEACCESSCONTROL
)

func (i BookingPageAccessControl) String() string {
    return []string{"unrestricted", "restrictedToOrganization", "unknownFutureValue"}[i]
}
func ParseBookingPageAccessControl(v string) (any, error) {
    result := UNRESTRICTED_BOOKINGPAGEACCESSCONTROL
    switch v {
        case "unrestricted":
            result = UNRESTRICTED_BOOKINGPAGEACCESSCONTROL
        case "restrictedToOrganization":
            result = RESTRICTEDTOORGANIZATION_BOOKINGPAGEACCESSCONTROL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGPAGEACCESSCONTROL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingPageAccessControl(values []BookingPageAccessControl) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingPageAccessControl) isMultiValue() bool {
    return false
}
