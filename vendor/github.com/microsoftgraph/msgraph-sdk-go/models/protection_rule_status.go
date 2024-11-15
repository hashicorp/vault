package models
type ProtectionRuleStatus int

const (
    DRAFT_PROTECTIONRULESTATUS ProtectionRuleStatus = iota
    ACTIVE_PROTECTIONRULESTATUS
    COMPLETED_PROTECTIONRULESTATUS
    COMPLETEDWITHERRORS_PROTECTIONRULESTATUS
    UNKNOWNFUTUREVALUE_PROTECTIONRULESTATUS
)

func (i ProtectionRuleStatus) String() string {
    return []string{"draft", "active", "completed", "completedWithErrors", "unknownFutureValue"}[i]
}
func ParseProtectionRuleStatus(v string) (any, error) {
    result := DRAFT_PROTECTIONRULESTATUS
    switch v {
        case "draft":
            result = DRAFT_PROTECTIONRULESTATUS
        case "active":
            result = ACTIVE_PROTECTIONRULESTATUS
        case "completed":
            result = COMPLETED_PROTECTIONRULESTATUS
        case "completedWithErrors":
            result = COMPLETEDWITHERRORS_PROTECTIONRULESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROTECTIONRULESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProtectionRuleStatus(values []ProtectionRuleStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProtectionRuleStatus) isMultiValue() bool {
    return false
}
