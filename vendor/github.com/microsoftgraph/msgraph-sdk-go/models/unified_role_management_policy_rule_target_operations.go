package models
type UnifiedRoleManagementPolicyRuleTargetOperations int

const (
    ALL_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS UnifiedRoleManagementPolicyRuleTargetOperations = iota
    ACTIVATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    DEACTIVATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    ASSIGN_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    UPDATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    REMOVE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    EXTEND_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    RENEW_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    UNKNOWNFUTUREVALUE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
)

func (i UnifiedRoleManagementPolicyRuleTargetOperations) String() string {
    return []string{"all", "activate", "deactivate", "assign", "update", "remove", "extend", "renew", "unknownFutureValue"}[i]
}
func ParseUnifiedRoleManagementPolicyRuleTargetOperations(v string) (any, error) {
    result := ALL_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
    switch v {
        case "all":
            result = ALL_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "activate":
            result = ACTIVATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "deactivate":
            result = DEACTIVATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "assign":
            result = ASSIGN_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "update":
            result = UPDATE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "remove":
            result = REMOVE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "extend":
            result = EXTEND_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "renew":
            result = RENEW_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_UNIFIEDROLEMANAGEMENTPOLICYRULETARGETOPERATIONS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUnifiedRoleManagementPolicyRuleTargetOperations(values []UnifiedRoleManagementPolicyRuleTargetOperations) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UnifiedRoleManagementPolicyRuleTargetOperations) isMultiValue() bool {
    return false
}
