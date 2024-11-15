package models
type IdentityUserFlowAttributeType int

const (
    BUILTIN_IDENTITYUSERFLOWATTRIBUTETYPE IdentityUserFlowAttributeType = iota
    CUSTOM_IDENTITYUSERFLOWATTRIBUTETYPE
    REQUIRED_IDENTITYUSERFLOWATTRIBUTETYPE
    UNKNOWNFUTUREVALUE_IDENTITYUSERFLOWATTRIBUTETYPE
)

func (i IdentityUserFlowAttributeType) String() string {
    return []string{"builtIn", "custom", "required", "unknownFutureValue"}[i]
}
func ParseIdentityUserFlowAttributeType(v string) (any, error) {
    result := BUILTIN_IDENTITYUSERFLOWATTRIBUTETYPE
    switch v {
        case "builtIn":
            result = BUILTIN_IDENTITYUSERFLOWATTRIBUTETYPE
        case "custom":
            result = CUSTOM_IDENTITYUSERFLOWATTRIBUTETYPE
        case "required":
            result = REQUIRED_IDENTITYUSERFLOWATTRIBUTETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_IDENTITYUSERFLOWATTRIBUTETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIdentityUserFlowAttributeType(values []IdentityUserFlowAttributeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IdentityUserFlowAttributeType) isMultiValue() bool {
    return false
}
