package models
type TrainingType int

const (
    UNKNOWN_TRAININGTYPE TrainingType = iota
    PHISHING_TRAININGTYPE
    UNKNOWNFUTUREVALUE_TRAININGTYPE
)

func (i TrainingType) String() string {
    return []string{"unknown", "phishing", "unknownFutureValue"}[i]
}
func ParseTrainingType(v string) (any, error) {
    result := UNKNOWN_TRAININGTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_TRAININGTYPE
        case "phishing":
            result = PHISHING_TRAININGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TRAININGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTrainingType(values []TrainingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TrainingType) isMultiValue() bool {
    return false
}
