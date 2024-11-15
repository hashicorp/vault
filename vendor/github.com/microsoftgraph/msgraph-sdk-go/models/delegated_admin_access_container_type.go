package models
type DelegatedAdminAccessContainerType int

const (
    SECURITYGROUP_DELEGATEDADMINACCESSCONTAINERTYPE DelegatedAdminAccessContainerType = iota
    UNKNOWNFUTUREVALUE_DELEGATEDADMINACCESSCONTAINERTYPE
)

func (i DelegatedAdminAccessContainerType) String() string {
    return []string{"securityGroup", "unknownFutureValue"}[i]
}
func ParseDelegatedAdminAccessContainerType(v string) (any, error) {
    result := SECURITYGROUP_DELEGATEDADMINACCESSCONTAINERTYPE
    switch v {
        case "securityGroup":
            result = SECURITYGROUP_DELEGATEDADMINACCESSCONTAINERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DELEGATEDADMINACCESSCONTAINERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDelegatedAdminAccessContainerType(values []DelegatedAdminAccessContainerType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DelegatedAdminAccessContainerType) isMultiValue() bool {
    return false
}
