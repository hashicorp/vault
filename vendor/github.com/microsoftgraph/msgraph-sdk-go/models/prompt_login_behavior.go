package models
type PromptLoginBehavior int

const (
    TRANSLATETOFRESHPASSWORDAUTHENTICATION_PROMPTLOGINBEHAVIOR PromptLoginBehavior = iota
    NATIVESUPPORT_PROMPTLOGINBEHAVIOR
    DISABLED_PROMPTLOGINBEHAVIOR
    UNKNOWNFUTUREVALUE_PROMPTLOGINBEHAVIOR
)

func (i PromptLoginBehavior) String() string {
    return []string{"translateToFreshPasswordAuthentication", "nativeSupport", "disabled", "unknownFutureValue"}[i]
}
func ParsePromptLoginBehavior(v string) (any, error) {
    result := TRANSLATETOFRESHPASSWORDAUTHENTICATION_PROMPTLOGINBEHAVIOR
    switch v {
        case "translateToFreshPasswordAuthentication":
            result = TRANSLATETOFRESHPASSWORDAUTHENTICATION_PROMPTLOGINBEHAVIOR
        case "nativeSupport":
            result = NATIVESUPPORT_PROMPTLOGINBEHAVIOR
        case "disabled":
            result = DISABLED_PROMPTLOGINBEHAVIOR
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROMPTLOGINBEHAVIOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePromptLoginBehavior(values []PromptLoginBehavior) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PromptLoginBehavior) isMultiValue() bool {
    return false
}
