package security
type IndicatorSource int

const (
    MICROSOFT_INDICATORSOURCE IndicatorSource = iota
    OSINT_INDICATORSOURCE
    PUBLIC_INDICATORSOURCE
    UNKNOWNFUTUREVALUE_INDICATORSOURCE
)

func (i IndicatorSource) String() string {
    return []string{"microsoft", "osint", "public", "unknownFutureValue"}[i]
}
func ParseIndicatorSource(v string) (any, error) {
    result := MICROSOFT_INDICATORSOURCE
    switch v {
        case "microsoft":
            result = MICROSOFT_INDICATORSOURCE
        case "osint":
            result = OSINT_INDICATORSOURCE
        case "public":
            result = PUBLIC_INDICATORSOURCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INDICATORSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIndicatorSource(values []IndicatorSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IndicatorSource) isMultiValue() bool {
    return false
}
