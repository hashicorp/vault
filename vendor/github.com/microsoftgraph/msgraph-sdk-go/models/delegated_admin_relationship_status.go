package models
type DelegatedAdminRelationshipStatus int

const (
    ACTIVATING_DELEGATEDADMINRELATIONSHIPSTATUS DelegatedAdminRelationshipStatus = iota
    ACTIVE_DELEGATEDADMINRELATIONSHIPSTATUS
    APPROVALPENDING_DELEGATEDADMINRELATIONSHIPSTATUS
    APPROVED_DELEGATEDADMINRELATIONSHIPSTATUS
    CREATED_DELEGATEDADMINRELATIONSHIPSTATUS
    EXPIRED_DELEGATEDADMINRELATIONSHIPSTATUS
    EXPIRING_DELEGATEDADMINRELATIONSHIPSTATUS
    TERMINATED_DELEGATEDADMINRELATIONSHIPSTATUS
    TERMINATING_DELEGATEDADMINRELATIONSHIPSTATUS
    TERMINATIONREQUESTED_DELEGATEDADMINRELATIONSHIPSTATUS
    UNKNOWNFUTUREVALUE_DELEGATEDADMINRELATIONSHIPSTATUS
)

func (i DelegatedAdminRelationshipStatus) String() string {
    return []string{"activating", "active", "approvalPending", "approved", "created", "expired", "expiring", "terminated", "terminating", "terminationRequested", "unknownFutureValue"}[i]
}
func ParseDelegatedAdminRelationshipStatus(v string) (any, error) {
    result := ACTIVATING_DELEGATEDADMINRELATIONSHIPSTATUS
    switch v {
        case "activating":
            result = ACTIVATING_DELEGATEDADMINRELATIONSHIPSTATUS
        case "active":
            result = ACTIVE_DELEGATEDADMINRELATIONSHIPSTATUS
        case "approvalPending":
            result = APPROVALPENDING_DELEGATEDADMINRELATIONSHIPSTATUS
        case "approved":
            result = APPROVED_DELEGATEDADMINRELATIONSHIPSTATUS
        case "created":
            result = CREATED_DELEGATEDADMINRELATIONSHIPSTATUS
        case "expired":
            result = EXPIRED_DELEGATEDADMINRELATIONSHIPSTATUS
        case "expiring":
            result = EXPIRING_DELEGATEDADMINRELATIONSHIPSTATUS
        case "terminated":
            result = TERMINATED_DELEGATEDADMINRELATIONSHIPSTATUS
        case "terminating":
            result = TERMINATING_DELEGATEDADMINRELATIONSHIPSTATUS
        case "terminationRequested":
            result = TERMINATIONREQUESTED_DELEGATEDADMINRELATIONSHIPSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DELEGATEDADMINRELATIONSHIPSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDelegatedAdminRelationshipStatus(values []DelegatedAdminRelationshipStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DelegatedAdminRelationshipStatus) isMultiValue() bool {
    return false
}
