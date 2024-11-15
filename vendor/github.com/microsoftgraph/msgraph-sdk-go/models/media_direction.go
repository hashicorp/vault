package models
type MediaDirection int

const (
    INACTIVE_MEDIADIRECTION MediaDirection = iota
    SENDONLY_MEDIADIRECTION
    RECEIVEONLY_MEDIADIRECTION
    SENDRECEIVE_MEDIADIRECTION
)

func (i MediaDirection) String() string {
    return []string{"inactive", "sendOnly", "receiveOnly", "sendReceive"}[i]
}
func ParseMediaDirection(v string) (any, error) {
    result := INACTIVE_MEDIADIRECTION
    switch v {
        case "inactive":
            result = INACTIVE_MEDIADIRECTION
        case "sendOnly":
            result = SENDONLY_MEDIADIRECTION
        case "receiveOnly":
            result = RECEIVEONLY_MEDIADIRECTION
        case "sendReceive":
            result = SENDRECEIVE_MEDIADIRECTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMediaDirection(values []MediaDirection) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MediaDirection) isMultiValue() bool {
    return false
}
