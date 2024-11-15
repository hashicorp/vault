package models
type ThreatAssessmentStatus int

const (
    PENDING_THREATASSESSMENTSTATUS ThreatAssessmentStatus = iota
    COMPLETED_THREATASSESSMENTSTATUS
)

func (i ThreatAssessmentStatus) String() string {
    return []string{"pending", "completed"}[i]
}
func ParseThreatAssessmentStatus(v string) (any, error) {
    result := PENDING_THREATASSESSMENTSTATUS
    switch v {
        case "pending":
            result = PENDING_THREATASSESSMENTSTATUS
        case "completed":
            result = COMPLETED_THREATASSESSMENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeThreatAssessmentStatus(values []ThreatAssessmentStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ThreatAssessmentStatus) isMultiValue() bool {
    return false
}
