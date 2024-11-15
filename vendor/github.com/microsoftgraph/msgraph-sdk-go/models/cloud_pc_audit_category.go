package models
type CloudPcAuditCategory int

const (
    CLOUDPC_CLOUDPCAUDITCATEGORY CloudPcAuditCategory = iota
    UNKNOWNFUTUREVALUE_CLOUDPCAUDITCATEGORY
)

func (i CloudPcAuditCategory) String() string {
    return []string{"cloudPC", "unknownFutureValue"}[i]
}
func ParseCloudPcAuditCategory(v string) (any, error) {
    result := CLOUDPC_CLOUDPCAUDITCATEGORY
    switch v {
        case "cloudPC":
            result = CLOUDPC_CLOUDPCAUDITCATEGORY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCAUDITCATEGORY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcAuditCategory(values []CloudPcAuditCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcAuditCategory) isMultiValue() bool {
    return false
}
