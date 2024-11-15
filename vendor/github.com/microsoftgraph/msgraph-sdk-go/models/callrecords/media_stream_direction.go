package callrecords
type MediaStreamDirection int

const (
    CALLERTOCALLEE_MEDIASTREAMDIRECTION MediaStreamDirection = iota
    CALLEETOCALLER_MEDIASTREAMDIRECTION
)

func (i MediaStreamDirection) String() string {
    return []string{"callerToCallee", "calleeToCaller"}[i]
}
func ParseMediaStreamDirection(v string) (any, error) {
    result := CALLERTOCALLEE_MEDIASTREAMDIRECTION
    switch v {
        case "callerToCallee":
            result = CALLERTOCALLEE_MEDIASTREAMDIRECTION
        case "calleeToCaller":
            result = CALLEETOCALLER_MEDIASTREAMDIRECTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMediaStreamDirection(values []MediaStreamDirection) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MediaStreamDirection) isMultiValue() bool {
    return false
}
