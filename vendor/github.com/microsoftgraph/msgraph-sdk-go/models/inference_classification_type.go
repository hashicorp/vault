package models
type InferenceClassificationType int

const (
    FOCUSED_INFERENCECLASSIFICATIONTYPE InferenceClassificationType = iota
    OTHER_INFERENCECLASSIFICATIONTYPE
)

func (i InferenceClassificationType) String() string {
    return []string{"focused", "other"}[i]
}
func ParseInferenceClassificationType(v string) (any, error) {
    result := FOCUSED_INFERENCECLASSIFICATIONTYPE
    switch v {
        case "focused":
            result = FOCUSED_INFERENCECLASSIFICATIONTYPE
        case "other":
            result = OTHER_INFERENCECLASSIFICATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeInferenceClassificationType(values []InferenceClassificationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i InferenceClassificationType) isMultiValue() bool {
    return false
}
