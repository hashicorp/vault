package models
type PrintQuality int

const (
    LOW_PRINTQUALITY PrintQuality = iota
    MEDIUM_PRINTQUALITY
    HIGH_PRINTQUALITY
    UNKNOWNFUTUREVALUE_PRINTQUALITY
)

func (i PrintQuality) String() string {
    return []string{"low", "medium", "high", "unknownFutureValue"}[i]
}
func ParsePrintQuality(v string) (any, error) {
    result := LOW_PRINTQUALITY
    switch v {
        case "low":
            result = LOW_PRINTQUALITY
        case "medium":
            result = MEDIUM_PRINTQUALITY
        case "high":
            result = HIGH_PRINTQUALITY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTQUALITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintQuality(values []PrintQuality) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintQuality) isMultiValue() bool {
    return false
}
