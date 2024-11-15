package models
type RegistryOperation int

const (
    UNKNOWN_REGISTRYOPERATION RegistryOperation = iota
    CREATE_REGISTRYOPERATION
    MODIFY_REGISTRYOPERATION
    DELETE_REGISTRYOPERATION
    UNKNOWNFUTUREVALUE_REGISTRYOPERATION
)

func (i RegistryOperation) String() string {
    return []string{"unknown", "create", "modify", "delete", "unknownFutureValue"}[i]
}
func ParseRegistryOperation(v string) (any, error) {
    result := UNKNOWN_REGISTRYOPERATION
    switch v {
        case "unknown":
            result = UNKNOWN_REGISTRYOPERATION
        case "create":
            result = CREATE_REGISTRYOPERATION
        case "modify":
            result = MODIFY_REGISTRYOPERATION
        case "delete":
            result = DELETE_REGISTRYOPERATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_REGISTRYOPERATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRegistryOperation(values []RegistryOperation) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RegistryOperation) isMultiValue() bool {
    return false
}
