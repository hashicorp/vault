package models
type DelegatedAdminAccessAssignmentStatus int

const (
    PENDING_DELEGATEDADMINACCESSASSIGNMENTSTATUS DelegatedAdminAccessAssignmentStatus = iota
    ACTIVE_DELEGATEDADMINACCESSASSIGNMENTSTATUS
    DELETING_DELEGATEDADMINACCESSASSIGNMENTSTATUS
    DELETED_DELEGATEDADMINACCESSASSIGNMENTSTATUS
    ERROR_DELEGATEDADMINACCESSASSIGNMENTSTATUS
    UNKNOWNFUTUREVALUE_DELEGATEDADMINACCESSASSIGNMENTSTATUS
)

func (i DelegatedAdminAccessAssignmentStatus) String() string {
    return []string{"pending", "active", "deleting", "deleted", "error", "unknownFutureValue"}[i]
}
func ParseDelegatedAdminAccessAssignmentStatus(v string) (any, error) {
    result := PENDING_DELEGATEDADMINACCESSASSIGNMENTSTATUS
    switch v {
        case "pending":
            result = PENDING_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        case "active":
            result = ACTIVE_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        case "deleting":
            result = DELETING_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        case "deleted":
            result = DELETED_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        case "error":
            result = ERROR_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DELEGATEDADMINACCESSASSIGNMENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDelegatedAdminAccessAssignmentStatus(values []DelegatedAdminAccessAssignmentStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DelegatedAdminAccessAssignmentStatus) isMultiValue() bool {
    return false
}
