package models
type SubjectRightsRequestType int

const (
    EXPORT_SUBJECTRIGHTSREQUESTTYPE SubjectRightsRequestType = iota
    DELETE_SUBJECTRIGHTSREQUESTTYPE
    ACCESS_SUBJECTRIGHTSREQUESTTYPE
    TAGFORACTION_SUBJECTRIGHTSREQUESTTYPE
    UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTTYPE
)

func (i SubjectRightsRequestType) String() string {
    return []string{"export", "delete", "access", "tagForAction", "unknownFutureValue"}[i]
}
func ParseSubjectRightsRequestType(v string) (any, error) {
    result := EXPORT_SUBJECTRIGHTSREQUESTTYPE
    switch v {
        case "export":
            result = EXPORT_SUBJECTRIGHTSREQUESTTYPE
        case "delete":
            result = DELETE_SUBJECTRIGHTSREQUESTTYPE
        case "access":
            result = ACCESS_SUBJECTRIGHTSREQUESTTYPE
        case "tagForAction":
            result = TAGFORACTION_SUBJECTRIGHTSREQUESTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SUBJECTRIGHTSREQUESTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSubjectRightsRequestType(values []SubjectRightsRequestType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SubjectRightsRequestType) isMultiValue() bool {
    return false
}
