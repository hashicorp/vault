package models
type UserType int

const (
    MEMBER_USERTYPE UserType = iota
    GUEST_USERTYPE
    UNKNOWNFUTUREVALUE_USERTYPE
)

func (i UserType) String() string {
    return []string{"member", "guest", "unknownFutureValue"}[i]
}
func ParseUserType(v string) (any, error) {
    result := MEMBER_USERTYPE
    switch v {
        case "member":
            result = MEMBER_USERTYPE
        case "guest":
            result = GUEST_USERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserType(values []UserType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserType) isMultiValue() bool {
    return false
}
