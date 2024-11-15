package models
type SharingDomainRestrictionMode int

const (
    NONE_SHARINGDOMAINRESTRICTIONMODE SharingDomainRestrictionMode = iota
    ALLOWLIST_SHARINGDOMAINRESTRICTIONMODE
    BLOCKLIST_SHARINGDOMAINRESTRICTIONMODE
    UNKNOWNFUTUREVALUE_SHARINGDOMAINRESTRICTIONMODE
)

func (i SharingDomainRestrictionMode) String() string {
    return []string{"none", "allowList", "blockList", "unknownFutureValue"}[i]
}
func ParseSharingDomainRestrictionMode(v string) (any, error) {
    result := NONE_SHARINGDOMAINRESTRICTIONMODE
    switch v {
        case "none":
            result = NONE_SHARINGDOMAINRESTRICTIONMODE
        case "allowList":
            result = ALLOWLIST_SHARINGDOMAINRESTRICTIONMODE
        case "blockList":
            result = BLOCKLIST_SHARINGDOMAINRESTRICTIONMODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SHARINGDOMAINRESTRICTIONMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSharingDomainRestrictionMode(values []SharingDomainRestrictionMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SharingDomainRestrictionMode) isMultiValue() bool {
    return false
}
