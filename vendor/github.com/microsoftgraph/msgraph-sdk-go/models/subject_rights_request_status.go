package models
type SubjectRightsRequestStatus int

const (
    ACTIVE_SUBJECTRIGHTSREQUESTSTATUS SubjectRightsRequestStatus = iota
    CLOSED_SUBJECTRIGHTSREQUESTSTATUS
    UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTSTATUS
)

func (i SubjectRightsRequestStatus) String() string {
    return []string{"active", "closed", "unknownFutureValue"}[i]
}
func ParseSubjectRightsRequestStatus(v string) (any, error) {
    result := ACTIVE_SUBJECTRIGHTSREQUESTSTATUS
    switch v {
        case "active":
            result = ACTIVE_SUBJECTRIGHTSREQUESTSTATUS
        case "closed":
            result = CLOSED_SUBJECTRIGHTSREQUESTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSubjectRightsRequestStatus(values []SubjectRightsRequestStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SubjectRightsRequestStatus) isMultiValue() bool {
    return false
}
