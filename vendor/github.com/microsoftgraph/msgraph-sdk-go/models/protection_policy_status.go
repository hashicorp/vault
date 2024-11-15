package models
type ProtectionPolicyStatus int

const (
    INACTIVE_PROTECTIONPOLICYSTATUS ProtectionPolicyStatus = iota
    ACTIVEWITHERRORS_PROTECTIONPOLICYSTATUS
    UPDATING_PROTECTIONPOLICYSTATUS
    ACTIVE_PROTECTIONPOLICYSTATUS
    UNKNOWNFUTUREVALUE_PROTECTIONPOLICYSTATUS
)

func (i ProtectionPolicyStatus) String() string {
    return []string{"inactive", "activeWithErrors", "updating", "active", "unknownFutureValue"}[i]
}
func ParseProtectionPolicyStatus(v string) (any, error) {
    result := INACTIVE_PROTECTIONPOLICYSTATUS
    switch v {
        case "inactive":
            result = INACTIVE_PROTECTIONPOLICYSTATUS
        case "activeWithErrors":
            result = ACTIVEWITHERRORS_PROTECTIONPOLICYSTATUS
        case "updating":
            result = UPDATING_PROTECTIONPOLICYSTATUS
        case "active":
            result = ACTIVE_PROTECTIONPOLICYSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROTECTIONPOLICYSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProtectionPolicyStatus(values []ProtectionPolicyStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProtectionPolicyStatus) isMultiValue() bool {
    return false
}
