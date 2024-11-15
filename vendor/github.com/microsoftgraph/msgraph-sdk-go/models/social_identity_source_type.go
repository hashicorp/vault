package models
type SocialIdentitySourceType int

const (
    FACEBOOK_SOCIALIDENTITYSOURCETYPE SocialIdentitySourceType = iota
    UNKNOWNFUTUREVALUE_SOCIALIDENTITYSOURCETYPE
)

func (i SocialIdentitySourceType) String() string {
    return []string{"facebook", "unknownFutureValue"}[i]
}
func ParseSocialIdentitySourceType(v string) (any, error) {
    result := FACEBOOK_SOCIALIDENTITYSOURCETYPE
    switch v {
        case "facebook":
            result = FACEBOOK_SOCIALIDENTITYSOURCETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SOCIALIDENTITYSOURCETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSocialIdentitySourceType(values []SocialIdentitySourceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SocialIdentitySourceType) isMultiValue() bool {
    return false
}
