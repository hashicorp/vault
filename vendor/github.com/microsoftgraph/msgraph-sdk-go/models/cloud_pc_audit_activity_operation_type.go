package models
type CloudPcAuditActivityOperationType int

const (
    CREATE_CLOUDPCAUDITACTIVITYOPERATIONTYPE CloudPcAuditActivityOperationType = iota
    DELETE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
    PATCH_CLOUDPCAUDITACTIVITYOPERATIONTYPE
    UNKNOWNFUTUREVALUE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
)

func (i CloudPcAuditActivityOperationType) String() string {
    return []string{"create", "delete", "patch", "unknownFutureValue"}[i]
}
func ParseCloudPcAuditActivityOperationType(v string) (any, error) {
    result := CREATE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
    switch v {
        case "create":
            result = CREATE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
        case "delete":
            result = DELETE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
        case "patch":
            result = PATCH_CLOUDPCAUDITACTIVITYOPERATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCAUDITACTIVITYOPERATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcAuditActivityOperationType(values []CloudPcAuditActivityOperationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcAuditActivityOperationType) isMultiValue() bool {
    return false
}
