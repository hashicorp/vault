package models
type DisableReason int

const (
    NONE_DISABLEREASON DisableReason = iota
    INVALIDBILLINGPROFILE_DISABLEREASON
    USERREQUESTED_DISABLEREASON
    UNKNOWNFUTUREVALUE_DISABLEREASON
)

func (i DisableReason) String() string {
    return []string{"none", "invalidBillingProfile", "userRequested", "unknownFutureValue"}[i]
}
func ParseDisableReason(v string) (any, error) {
    result := NONE_DISABLEREASON
    switch v {
        case "none":
            result = NONE_DISABLEREASON
        case "invalidBillingProfile":
            result = INVALIDBILLINGPROFILE_DISABLEREASON
        case "userRequested":
            result = USERREQUESTED_DISABLEREASON
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DISABLEREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDisableReason(values []DisableReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DisableReason) isMultiValue() bool {
    return false
}
