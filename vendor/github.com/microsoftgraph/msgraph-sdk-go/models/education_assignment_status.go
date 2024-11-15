package models
type EducationAssignmentStatus int

const (
    DRAFT_EDUCATIONASSIGNMENTSTATUS EducationAssignmentStatus = iota
    PUBLISHED_EDUCATIONASSIGNMENTSTATUS
    ASSIGNED_EDUCATIONASSIGNMENTSTATUS
    UNKNOWNFUTUREVALUE_EDUCATIONASSIGNMENTSTATUS
    INACTIVE_EDUCATIONASSIGNMENTSTATUS
)

func (i EducationAssignmentStatus) String() string {
    return []string{"draft", "published", "assigned", "unknownFutureValue", "inactive"}[i]
}
func ParseEducationAssignmentStatus(v string) (any, error) {
    result := DRAFT_EDUCATIONASSIGNMENTSTATUS
    switch v {
        case "draft":
            result = DRAFT_EDUCATIONASSIGNMENTSTATUS
        case "published":
            result = PUBLISHED_EDUCATIONASSIGNMENTSTATUS
        case "assigned":
            result = ASSIGNED_EDUCATIONASSIGNMENTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EDUCATIONASSIGNMENTSTATUS
        case "inactive":
            result = INACTIVE_EDUCATIONASSIGNMENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEducationAssignmentStatus(values []EducationAssignmentStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EducationAssignmentStatus) isMultiValue() bool {
    return false
}
