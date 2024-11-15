package models
type UserSignInRecommendationScope int

const (
    TENANT_USERSIGNINRECOMMENDATIONSCOPE UserSignInRecommendationScope = iota
    APPLICATION_USERSIGNINRECOMMENDATIONSCOPE
    UNKNOWNFUTUREVALUE_USERSIGNINRECOMMENDATIONSCOPE
)

func (i UserSignInRecommendationScope) String() string {
    return []string{"tenant", "application", "unknownFutureValue"}[i]
}
func ParseUserSignInRecommendationScope(v string) (any, error) {
    result := TENANT_USERSIGNINRECOMMENDATIONSCOPE
    switch v {
        case "tenant":
            result = TENANT_USERSIGNINRECOMMENDATIONSCOPE
        case "application":
            result = APPLICATION_USERSIGNINRECOMMENDATIONSCOPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USERSIGNINRECOMMENDATIONSCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserSignInRecommendationScope(values []UserSignInRecommendationScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserSignInRecommendationScope) isMultiValue() bool {
    return false
}
