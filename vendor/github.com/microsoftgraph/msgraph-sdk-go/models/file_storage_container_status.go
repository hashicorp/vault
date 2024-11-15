package models
type FileStorageContainerStatus int

const (
    INACTIVE_FILESTORAGECONTAINERSTATUS FileStorageContainerStatus = iota
    ACTIVE_FILESTORAGECONTAINERSTATUS
    UNKNOWNFUTUREVALUE_FILESTORAGECONTAINERSTATUS
)

func (i FileStorageContainerStatus) String() string {
    return []string{"inactive", "active", "unknownFutureValue"}[i]
}
func ParseFileStorageContainerStatus(v string) (any, error) {
    result := INACTIVE_FILESTORAGECONTAINERSTATUS
    switch v {
        case "inactive":
            result = INACTIVE_FILESTORAGECONTAINERSTATUS
        case "active":
            result = ACTIVE_FILESTORAGECONTAINERSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FILESTORAGECONTAINERSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFileStorageContainerStatus(values []FileStorageContainerStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FileStorageContainerStatus) isMultiValue() bool {
    return false
}
