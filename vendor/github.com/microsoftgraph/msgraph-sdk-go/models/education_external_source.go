package models
type EducationExternalSource int

const (
    SIS_EDUCATIONEXTERNALSOURCE EducationExternalSource = iota
    MANUAL_EDUCATIONEXTERNALSOURCE
    UNKNOWNFUTUREVALUE_EDUCATIONEXTERNALSOURCE
)

func (i EducationExternalSource) String() string {
    return []string{"sis", "manual", "unknownFutureValue"}[i]
}
func ParseEducationExternalSource(v string) (any, error) {
    result := SIS_EDUCATIONEXTERNALSOURCE
    switch v {
        case "sis":
            result = SIS_EDUCATIONEXTERNALSOURCE
        case "manual":
            result = MANUAL_EDUCATIONEXTERNALSOURCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EDUCATIONEXTERNALSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEducationExternalSource(values []EducationExternalSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EducationExternalSource) isMultiValue() bool {
    return false
}
