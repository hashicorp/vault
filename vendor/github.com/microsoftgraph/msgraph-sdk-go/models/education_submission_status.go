package models
type EducationSubmissionStatus int

const (
    WORKING_EDUCATIONSUBMISSIONSTATUS EducationSubmissionStatus = iota
    SUBMITTED_EDUCATIONSUBMISSIONSTATUS
    RELEASED_EDUCATIONSUBMISSIONSTATUS
    RETURNED_EDUCATIONSUBMISSIONSTATUS
    UNKNOWNFUTUREVALUE_EDUCATIONSUBMISSIONSTATUS
    REASSIGNED_EDUCATIONSUBMISSIONSTATUS
    EXCUSED_EDUCATIONSUBMISSIONSTATUS
)

func (i EducationSubmissionStatus) String() string {
    return []string{"working", "submitted", "released", "returned", "unknownFutureValue", "reassigned", "excused"}[i]
}
func ParseEducationSubmissionStatus(v string) (any, error) {
    result := WORKING_EDUCATIONSUBMISSIONSTATUS
    switch v {
        case "working":
            result = WORKING_EDUCATIONSUBMISSIONSTATUS
        case "submitted":
            result = SUBMITTED_EDUCATIONSUBMISSIONSTATUS
        case "released":
            result = RELEASED_EDUCATIONSUBMISSIONSTATUS
        case "returned":
            result = RETURNED_EDUCATIONSUBMISSIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EDUCATIONSUBMISSIONSTATUS
        case "reassigned":
            result = REASSIGNED_EDUCATIONSUBMISSIONSTATUS
        case "excused":
            result = EXCUSED_EDUCATIONSUBMISSIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEducationSubmissionStatus(values []EducationSubmissionStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EducationSubmissionStatus) isMultiValue() bool {
    return false
}
