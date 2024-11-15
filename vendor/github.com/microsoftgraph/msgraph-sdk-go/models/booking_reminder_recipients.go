package models
type BookingReminderRecipients int

const (
    ALLATTENDEES_BOOKINGREMINDERRECIPIENTS BookingReminderRecipients = iota
    STAFF_BOOKINGREMINDERRECIPIENTS
    CUSTOMER_BOOKINGREMINDERRECIPIENTS
    UNKNOWNFUTUREVALUE_BOOKINGREMINDERRECIPIENTS
)

func (i BookingReminderRecipients) String() string {
    return []string{"allAttendees", "staff", "customer", "unknownFutureValue"}[i]
}
func ParseBookingReminderRecipients(v string) (any, error) {
    result := ALLATTENDEES_BOOKINGREMINDERRECIPIENTS
    switch v {
        case "allAttendees":
            result = ALLATTENDEES_BOOKINGREMINDERRECIPIENTS
        case "staff":
            result = STAFF_BOOKINGREMINDERRECIPIENTS
        case "customer":
            result = CUSTOMER_BOOKINGREMINDERRECIPIENTS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BOOKINGREMINDERRECIPIENTS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBookingReminderRecipients(values []BookingReminderRecipients) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BookingReminderRecipients) isMultiValue() bool {
    return false
}
