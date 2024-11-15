package models
type ResponseType int

const (
    NONE_RESPONSETYPE ResponseType = iota
    ORGANIZER_RESPONSETYPE
    TENTATIVELYACCEPTED_RESPONSETYPE
    ACCEPTED_RESPONSETYPE
    DECLINED_RESPONSETYPE
    NOTRESPONDED_RESPONSETYPE
)

func (i ResponseType) String() string {
    return []string{"none", "organizer", "tentativelyAccepted", "accepted", "declined", "notResponded"}[i]
}
func ParseResponseType(v string) (any, error) {
    result := NONE_RESPONSETYPE
    switch v {
        case "none":
            result = NONE_RESPONSETYPE
        case "organizer":
            result = ORGANIZER_RESPONSETYPE
        case "tentativelyAccepted":
            result = TENTATIVELYACCEPTED_RESPONSETYPE
        case "accepted":
            result = ACCEPTED_RESPONSETYPE
        case "declined":
            result = DECLINED_RESPONSETYPE
        case "notResponded":
            result = NOTRESPONDED_RESPONSETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeResponseType(values []ResponseType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ResponseType) isMultiValue() bool {
    return false
}
