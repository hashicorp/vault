package models
type BookingStaffMembershipStatus int

const (
    ACTIVE_BOOKINGSTAFFMEMBERSHIPSTATUS BookingStaffMembershipStatus = iota
    PENDINGACCEPTANCE_BOOKINGSTAFFMEMBERSHIPSTATUS
    REJECTEDBYSTAFF_BOOKINGSTAFFMEMBERSHIPSTATUS
    UNKNOWNFUTUREVALUE_BOOKINGSTAFFMEMBERSHIPSTATUS
)

func (i BookingStaffMembershipStatus) String() string {
    return []string{"active", "pendingAcceptance", "rejectedByStaff", "unknownFutureValue"}[i]
}
func ParseBookingStaffMembershipStatus(v string) (any, error) {
    result := ACTIVE_BOOKINGSTAFFMEMBERSHIPSTATUS
    switch v {
        case "active":
            result = ACTIVE_BOOKINGSTAFFMEMBERSHIPSTATUS
        case "pendingAcceptance":
            result = PENDINGACCEPTANCE_BOOKINGSTAFFMEMBERSHIPSTATUS
        case "rejectedByStaff":
            result = REJECTEDBYSTAFF_BOOKINGSTAFFMEMBERSHIPSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGSTAFFMEMBERSHIPSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingStaffMembershipStatus(values []BookingStaffMembershipStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingStaffMembershipStatus) isMultiValue() bool {
    return false
}
