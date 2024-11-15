package models
// Android Work Profile default app permission policy type.
type AndroidWorkProfileDefaultAppPermissionPolicyType int

const (
    // Device default value, no intent.
    DEVICEDEFAULT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE AndroidWorkProfileDefaultAppPermissionPolicyType = iota
    // Prompt.
    PROMPT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
    // Auto grant.
    AUTOGRANT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
    // Auto deny.
    AUTODENY_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
)

func (i AndroidWorkProfileDefaultAppPermissionPolicyType) String() string {
    return []string{"deviceDefault", "prompt", "autoGrant", "autoDeny"}[i]
}
func ParseAndroidWorkProfileDefaultAppPermissionPolicyType(v string) (any, error) {
    result := DEVICEDEFAULT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
    switch v {
        case "deviceDefault":
            result = DEVICEDEFAULT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
        case "prompt":
            result = PROMPT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
        case "autoGrant":
            result = AUTOGRANT_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
        case "autoDeny":
            result = AUTODENY_ANDROIDWORKPROFILEDEFAULTAPPPERMISSIONPOLICYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAndroidWorkProfileDefaultAppPermissionPolicyType(values []AndroidWorkProfileDefaultAppPermissionPolicyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AndroidWorkProfileDefaultAppPermissionPolicyType) isMultiValue() bool {
    return false
}
