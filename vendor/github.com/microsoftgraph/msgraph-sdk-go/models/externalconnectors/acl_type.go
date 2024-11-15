package externalconnectors
type AclType int

const (
    USER_ACLTYPE AclType = iota
    GROUP_ACLTYPE
    EVERYONE_ACLTYPE
    EVERYONEEXCEPTGUESTS_ACLTYPE
    EXTERNALGROUP_ACLTYPE
    UNKNOWNFUTUREVALUE_ACLTYPE
)

func (i AclType) String() string {
    return []string{"user", "group", "everyone", "everyoneExceptGuests", "externalGroup", "unknownFutureValue"}[i]
}
func ParseAclType(v string) (any, error) {
    result := USER_ACLTYPE
    switch v {
        case "user":
            result = USER_ACLTYPE
        case "group":
            result = GROUP_ACLTYPE
        case "everyone":
            result = EVERYONE_ACLTYPE
        case "everyoneExceptGuests":
            result = EVERYONEEXCEPTGUESTS_ACLTYPE
        case "externalGroup":
            result = EXTERNALGROUP_ACLTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACLTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAclType(values []AclType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AclType) isMultiValue() bool {
    return false
}
