package models
type PermissionClassificationType int

const (
    LOW_PERMISSIONCLASSIFICATIONTYPE PermissionClassificationType = iota
    MEDIUM_PERMISSIONCLASSIFICATIONTYPE
    HIGH_PERMISSIONCLASSIFICATIONTYPE
    UNKNOWNFUTUREVALUE_PERMISSIONCLASSIFICATIONTYPE
)

func (i PermissionClassificationType) String() string {
    return []string{"low", "medium", "high", "unknownFutureValue"}[i]
}
func ParsePermissionClassificationType(v string) (any, error) {
    result := LOW_PERMISSIONCLASSIFICATIONTYPE
    switch v {
        case "low":
            result = LOW_PERMISSIONCLASSIFICATIONTYPE
        case "medium":
            result = MEDIUM_PERMISSIONCLASSIFICATIONTYPE
        case "high":
            result = HIGH_PERMISSIONCLASSIFICATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PERMISSIONCLASSIFICATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePermissionClassificationType(values []PermissionClassificationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PermissionClassificationType) isMultiValue() bool {
    return false
}
