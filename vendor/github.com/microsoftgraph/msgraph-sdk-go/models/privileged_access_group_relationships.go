package models
type PrivilegedAccessGroupRelationships int

const (
    OWNER_PRIVILEGEDACCESSGROUPRELATIONSHIPS PrivilegedAccessGroupRelationships = iota
    MEMBER_PRIVILEGEDACCESSGROUPRELATIONSHIPS
    UNKNOWNFUTUREVALUE_PRIVILEGEDACCESSGROUPRELATIONSHIPS
)

func (i PrivilegedAccessGroupRelationships) String() string {
    return []string{"owner", "member", "unknownFutureValue"}[i]
}
func ParsePrivilegedAccessGroupRelationships(v string) (any, error) {
    result := OWNER_PRIVILEGEDACCESSGROUPRELATIONSHIPS
    switch v {
        case "owner":
            result = OWNER_PRIVILEGEDACCESSGROUPRELATIONSHIPS
        case "member":
            result = MEMBER_PRIVILEGEDACCESSGROUPRELATIONSHIPS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRIVILEGEDACCESSGROUPRELATIONSHIPS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrivilegedAccessGroupRelationships(values []PrivilegedAccessGroupRelationships) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrivilegedAccessGroupRelationships) isMultiValue() bool {
    return false
}
