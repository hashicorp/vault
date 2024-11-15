package models
type TargettedUserType int

const (
    UNKNOWN_TARGETTEDUSERTYPE TargettedUserType = iota
    CLICKED_TARGETTEDUSERTYPE
    COMPROMISED_TARGETTEDUSERTYPE
    ALLUSERS_TARGETTEDUSERTYPE
    UNKNOWNFUTUREVALUE_TARGETTEDUSERTYPE
)

func (i TargettedUserType) String() string {
    return []string{"unknown", "clicked", "compromised", "allUsers", "unknownFutureValue"}[i]
}
func ParseTargettedUserType(v string) (any, error) {
    result := UNKNOWN_TARGETTEDUSERTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_TARGETTEDUSERTYPE
        case "clicked":
            result = CLICKED_TARGETTEDUSERTYPE
        case "compromised":
            result = COMPROMISED_TARGETTEDUSERTYPE
        case "allUsers":
            result = ALLUSERS_TARGETTEDUSERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TARGETTEDUSERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTargettedUserType(values []TargettedUserType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TargettedUserType) isMultiValue() bool {
    return false
}
