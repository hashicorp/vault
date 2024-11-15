package models
type ThreatAssessmentRequestSource int

const (
    UNDEFINED_THREATASSESSMENTREQUESTSOURCE ThreatAssessmentRequestSource = iota
    USER_THREATASSESSMENTREQUESTSOURCE
    ADMINISTRATOR_THREATASSESSMENTREQUESTSOURCE
)

func (i ThreatAssessmentRequestSource) String() string {
    return []string{"undefined", "user", "administrator"}[i]
}
func ParseThreatAssessmentRequestSource(v string) (any, error) {
    result := UNDEFINED_THREATASSESSMENTREQUESTSOURCE
    switch v {
        case "undefined":
            result = UNDEFINED_THREATASSESSMENTREQUESTSOURCE
        case "user":
            result = USER_THREATASSESSMENTREQUESTSOURCE
        case "administrator":
            result = ADMINISTRATOR_THREATASSESSMENTREQUESTSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeThreatAssessmentRequestSource(values []ThreatAssessmentRequestSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ThreatAssessmentRequestSource) isMultiValue() bool {
    return false
}
