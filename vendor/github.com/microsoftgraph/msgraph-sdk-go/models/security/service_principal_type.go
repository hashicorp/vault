package security
type ServicePrincipalType int

const (
    UNKNOWN_SERVICEPRINCIPALTYPE ServicePrincipalType = iota
    APPLICATION_SERVICEPRINCIPALTYPE
    MANAGEDIDENTITY_SERVICEPRINCIPALTYPE
    LEGACY_SERVICEPRINCIPALTYPE
    UNKNOWNFUTUREVALUE_SERVICEPRINCIPALTYPE
)

func (i ServicePrincipalType) String() string {
    return []string{"unknown", "application", "managedIdentity", "legacy", "unknownFutureValue"}[i]
}
func ParseServicePrincipalType(v string) (any, error) {
    result := UNKNOWN_SERVICEPRINCIPALTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_SERVICEPRINCIPALTYPE
        case "application":
            result = APPLICATION_SERVICEPRINCIPALTYPE
        case "managedIdentity":
            result = MANAGEDIDENTITY_SERVICEPRINCIPALTYPE
        case "legacy":
            result = LEGACY_SERVICEPRINCIPALTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICEPRINCIPALTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServicePrincipalType(values []ServicePrincipalType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServicePrincipalType) isMultiValue() bool {
    return false
}
