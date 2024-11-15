package models
// Possible values to specify which cookies are allowed in Microsoft Edge.
type EdgeCookiePolicy int

const (
    // Allow the user to set.
    USERDEFINED_EDGECOOKIEPOLICY EdgeCookiePolicy = iota
    // Allow.
    ALLOW_EDGECOOKIEPOLICY
    // Block only third party cookies.
    BLOCKTHIRDPARTY_EDGECOOKIEPOLICY
    // Block all cookies.
    BLOCKALL_EDGECOOKIEPOLICY
)

func (i EdgeCookiePolicy) String() string {
    return []string{"userDefined", "allow", "blockThirdParty", "blockAll"}[i]
}
func ParseEdgeCookiePolicy(v string) (any, error) {
    result := USERDEFINED_EDGECOOKIEPOLICY
    switch v {
        case "userDefined":
            result = USERDEFINED_EDGECOOKIEPOLICY
        case "allow":
            result = ALLOW_EDGECOOKIEPOLICY
        case "blockThirdParty":
            result = BLOCKTHIRDPARTY_EDGECOOKIEPOLICY
        case "blockAll":
            result = BLOCKALL_EDGECOOKIEPOLICY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEdgeCookiePolicy(values []EdgeCookiePolicy) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EdgeCookiePolicy) isMultiValue() bool {
    return false
}
