package models
type SignInUserType int

const (
    MEMBER_SIGNINUSERTYPE SignInUserType = iota
    GUEST_SIGNINUSERTYPE
    UNKNOWNFUTUREVALUE_SIGNINUSERTYPE
)

func (i SignInUserType) String() string {
    return []string{"member", "guest", "unknownFutureValue"}[i]
}
func ParseSignInUserType(v string) (any, error) {
    result := MEMBER_SIGNINUSERTYPE
    switch v {
        case "member":
            result = MEMBER_SIGNINUSERTYPE
        case "guest":
            result = GUEST_SIGNINUSERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIGNINUSERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSignInUserType(values []SignInUserType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SignInUserType) isMultiValue() bool {
    return false
}
