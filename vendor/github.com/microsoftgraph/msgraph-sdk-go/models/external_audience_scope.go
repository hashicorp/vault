package models
type ExternalAudienceScope int

const (
    NONE_EXTERNALAUDIENCESCOPE ExternalAudienceScope = iota
    CONTACTSONLY_EXTERNALAUDIENCESCOPE
    ALL_EXTERNALAUDIENCESCOPE
)

func (i ExternalAudienceScope) String() string {
    return []string{"none", "contactsOnly", "all"}[i]
}
func ParseExternalAudienceScope(v string) (any, error) {
    result := NONE_EXTERNALAUDIENCESCOPE
    switch v {
        case "none":
            result = NONE_EXTERNALAUDIENCESCOPE
        case "contactsOnly":
            result = CONTACTSONLY_EXTERNALAUDIENCESCOPE
        case "all":
            result = ALL_EXTERNALAUDIENCESCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExternalAudienceScope(values []ExternalAudienceScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExternalAudienceScope) isMultiValue() bool {
    return false
}
