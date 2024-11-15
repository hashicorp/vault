package models
type DelegatedAdminRelationshipOperationType int

const (
    DELEGATEDADMINACCESSASSIGNMENTUPDATE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE DelegatedAdminRelationshipOperationType = iota
    UNKNOWNFUTUREVALUE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
    DELEGATEDADMINRELATIONSHIPUPDATE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
)

func (i DelegatedAdminRelationshipOperationType) String() string {
    return []string{"delegatedAdminAccessAssignmentUpdate", "unknownFutureValue", "delegatedAdminRelationshipUpdate"}[i]
}
func ParseDelegatedAdminRelationshipOperationType(v string) (any, error) {
    result := DELEGATEDADMINACCESSASSIGNMENTUPDATE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
    switch v {
        case "delegatedAdminAccessAssignmentUpdate":
            result = DELEGATEDADMINACCESSASSIGNMENTUPDATE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
        case "delegatedAdminRelationshipUpdate":
            result = DELEGATEDADMINRELATIONSHIPUPDATE_DELEGATEDADMINRELATIONSHIPOPERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDelegatedAdminRelationshipOperationType(values []DelegatedAdminRelationshipOperationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DelegatedAdminRelationshipOperationType) isMultiValue() bool {
    return false
}
