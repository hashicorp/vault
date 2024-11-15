package models
type ProcessIntegrityLevel int

const (
    UNKNOWN_PROCESSINTEGRITYLEVEL ProcessIntegrityLevel = iota
    UNTRUSTED_PROCESSINTEGRITYLEVEL
    LOW_PROCESSINTEGRITYLEVEL
    MEDIUM_PROCESSINTEGRITYLEVEL
    HIGH_PROCESSINTEGRITYLEVEL
    SYSTEM_PROCESSINTEGRITYLEVEL
    UNKNOWNFUTUREVALUE_PROCESSINTEGRITYLEVEL
)

func (i ProcessIntegrityLevel) String() string {
    return []string{"unknown", "untrusted", "low", "medium", "high", "system", "unknownFutureValue"}[i]
}
func ParseProcessIntegrityLevel(v string) (any, error) {
    result := UNKNOWN_PROCESSINTEGRITYLEVEL
    switch v {
        case "unknown":
            result = UNKNOWN_PROCESSINTEGRITYLEVEL
        case "untrusted":
            result = UNTRUSTED_PROCESSINTEGRITYLEVEL
        case "low":
            result = LOW_PROCESSINTEGRITYLEVEL
        case "medium":
            result = MEDIUM_PROCESSINTEGRITYLEVEL
        case "high":
            result = HIGH_PROCESSINTEGRITYLEVEL
        case "system":
            result = SYSTEM_PROCESSINTEGRITYLEVEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROCESSINTEGRITYLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProcessIntegrityLevel(values []ProcessIntegrityLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProcessIntegrityLevel) isMultiValue() bool {
    return false
}
