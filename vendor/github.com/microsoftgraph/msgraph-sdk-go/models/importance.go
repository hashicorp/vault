package models
type Importance int

const (
    LOW_IMPORTANCE Importance = iota
    NORMAL_IMPORTANCE
    HIGH_IMPORTANCE
)

func (i Importance) String() string {
    return []string{"low", "normal", "high"}[i]
}
func ParseImportance(v string) (any, error) {
    result := LOW_IMPORTANCE
    switch v {
        case "low":
            result = LOW_IMPORTANCE
        case "normal":
            result = NORMAL_IMPORTANCE
        case "high":
            result = HIGH_IMPORTANCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeImportance(values []Importance) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Importance) isMultiValue() bool {
    return false
}
