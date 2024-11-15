package models
type ThreatAssessmentContentType int

const (
    MAIL_THREATASSESSMENTCONTENTTYPE ThreatAssessmentContentType = iota
    URL_THREATASSESSMENTCONTENTTYPE
    FILE_THREATASSESSMENTCONTENTTYPE
)

func (i ThreatAssessmentContentType) String() string {
    return []string{"mail", "url", "file"}[i]
}
func ParseThreatAssessmentContentType(v string) (any, error) {
    result := MAIL_THREATASSESSMENTCONTENTTYPE
    switch v {
        case "mail":
            result = MAIL_THREATASSESSMENTCONTENTTYPE
        case "url":
            result = URL_THREATASSESSMENTCONTENTTYPE
        case "file":
            result = FILE_THREATASSESSMENTCONTENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeThreatAssessmentContentType(values []ThreatAssessmentContentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ThreatAssessmentContentType) isMultiValue() bool {
    return false
}
