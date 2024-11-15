package models
type EnrollmentState int

const (
    // Device enrollment state is unknown
    UNKNOWN_ENROLLMENTSTATE EnrollmentState = iota
    // Device is Enrolled.
    ENROLLED_ENROLLMENTSTATE
    // Enrolled but it's enrolled via enrollment profile and the enrolled profile is different from the assigned profile.
    PENDINGRESET_ENROLLMENTSTATE
    // Not enrolled and there is enrollment failure record.
    FAILED_ENROLLMENTSTATE
    // Device is imported but not enrolled.
    NOTCONTACTED_ENROLLMENTSTATE
)

func (i EnrollmentState) String() string {
    return []string{"unknown", "enrolled", "pendingReset", "failed", "notContacted"}[i]
}
func ParseEnrollmentState(v string) (any, error) {
    result := UNKNOWN_ENROLLMENTSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_ENROLLMENTSTATE
        case "enrolled":
            result = ENROLLED_ENROLLMENTSTATE
        case "pendingReset":
            result = PENDINGRESET_ENROLLMENTSTATE
        case "failed":
            result = FAILED_ENROLLMENTSTATE
        case "notContacted":
            result = NOTCONTACTED_ENROLLMENTSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEnrollmentState(values []EnrollmentState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EnrollmentState) isMultiValue() bool {
    return false
}
