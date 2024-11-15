package models
type EducationModuleStatus int

const (
    DRAFT_EDUCATIONMODULESTATUS EducationModuleStatus = iota
    PUBLISHED_EDUCATIONMODULESTATUS
    UNKNOWNFUTUREVALUE_EDUCATIONMODULESTATUS
)

func (i EducationModuleStatus) String() string {
    return []string{"draft", "published", "unknownFutureValue"}[i]
}
func ParseEducationModuleStatus(v string) (any, error) {
    result := DRAFT_EDUCATIONMODULESTATUS
    switch v {
        case "draft":
            result = DRAFT_EDUCATIONMODULESTATUS
        case "published":
            result = PUBLISHED_EDUCATIONMODULESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EDUCATIONMODULESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEducationModuleStatus(values []EducationModuleStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EducationModuleStatus) isMultiValue() bool {
    return false
}
