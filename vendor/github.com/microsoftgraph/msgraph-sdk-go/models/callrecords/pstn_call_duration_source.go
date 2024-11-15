package callrecords
type PstnCallDurationSource int

const (
    MICROSOFT_PSTNCALLDURATIONSOURCE PstnCallDurationSource = iota
    OPERATOR_PSTNCALLDURATIONSOURCE
)

func (i PstnCallDurationSource) String() string {
    return []string{"microsoft", "operator"}[i]
}
func ParsePstnCallDurationSource(v string) (any, error) {
    result := MICROSOFT_PSTNCALLDURATIONSOURCE
    switch v {
        case "microsoft":
            result = MICROSOFT_PSTNCALLDURATIONSOURCE
        case "operator":
            result = OPERATOR_PSTNCALLDURATIONSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePstnCallDurationSource(values []PstnCallDurationSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PstnCallDurationSource) isMultiValue() bool {
    return false
}
