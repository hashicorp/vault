package models
type SubjectRightsRequestStageStatus int

const (
    NOTSTARTED_SUBJECTRIGHTSREQUESTSTAGESTATUS SubjectRightsRequestStageStatus = iota
    CURRENT_SUBJECTRIGHTSREQUESTSTAGESTATUS
    COMPLETED_SUBJECTRIGHTSREQUESTSTAGESTATUS
    FAILED_SUBJECTRIGHTSREQUESTSTAGESTATUS
    UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTSTAGESTATUS
)

func (i SubjectRightsRequestStageStatus) String() string {
    return []string{"notStarted", "current", "completed", "failed", "unknownFutureValue"}[i]
}
func ParseSubjectRightsRequestStageStatus(v string) (any, error) {
    result := NOTSTARTED_SUBJECTRIGHTSREQUESTSTAGESTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_SUBJECTRIGHTSREQUESTSTAGESTATUS
        case "current":
            result = CURRENT_SUBJECTRIGHTSREQUESTSTAGESTATUS
        case "completed":
            result = COMPLETED_SUBJECTRIGHTSREQUESTSTAGESTATUS
        case "failed":
            result = FAILED_SUBJECTRIGHTSREQUESTSTAGESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTSTAGESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSubjectRightsRequestStageStatus(values []SubjectRightsRequestStageStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SubjectRightsRequestStageStatus) isMultiValue() bool {
    return false
}
