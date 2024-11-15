package models
type ServiceAppStatus int

const (
    INACTIVE_SERVICEAPPSTATUS ServiceAppStatus = iota
    ACTIVE_SERVICEAPPSTATUS
    PENDINGACTIVE_SERVICEAPPSTATUS
    PENDINGINACTIVE_SERVICEAPPSTATUS
    UNKNOWNFUTUREVALUE_SERVICEAPPSTATUS
)

func (i ServiceAppStatus) String() string {
    return []string{"inactive", "active", "pendingActive", "pendingInactive", "unknownFutureValue"}[i]
}
func ParseServiceAppStatus(v string) (any, error) {
    result := INACTIVE_SERVICEAPPSTATUS
    switch v {
        case "inactive":
            result = INACTIVE_SERVICEAPPSTATUS
        case "active":
            result = ACTIVE_SERVICEAPPSTATUS
        case "pendingActive":
            result = PENDINGACTIVE_SERVICEAPPSTATUS
        case "pendingInactive":
            result = PENDINGINACTIVE_SERVICEAPPSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICEAPPSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServiceAppStatus(values []ServiceAppStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServiceAppStatus) isMultiValue() bool {
    return false
}
