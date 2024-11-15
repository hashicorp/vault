package identitygovernance
type MembershipChangeType int

const (
    ADD_MEMBERSHIPCHANGETYPE MembershipChangeType = iota
    REMOVE_MEMBERSHIPCHANGETYPE
    UNKNOWNFUTUREVALUE_MEMBERSHIPCHANGETYPE
)

func (i MembershipChangeType) String() string {
    return []string{"add", "remove", "unknownFutureValue"}[i]
}
func ParseMembershipChangeType(v string) (any, error) {
    result := ADD_MEMBERSHIPCHANGETYPE
    switch v {
        case "add":
            result = ADD_MEMBERSHIPCHANGETYPE
        case "remove":
            result = REMOVE_MEMBERSHIPCHANGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MEMBERSHIPCHANGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMembershipChangeType(values []MembershipChangeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MembershipChangeType) isMultiValue() bool {
    return false
}
