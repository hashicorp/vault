package security
type HostPortStatus int

const (
    OPEN_HOSTPORTSTATUS HostPortStatus = iota
    FILTERED_HOSTPORTSTATUS
    CLOSED_HOSTPORTSTATUS
    UNKNOWNFUTUREVALUE_HOSTPORTSTATUS
)

func (i HostPortStatus) String() string {
    return []string{"open", "filtered", "closed", "unknownFutureValue"}[i]
}
func ParseHostPortStatus(v string) (any, error) {
    result := OPEN_HOSTPORTSTATUS
    switch v {
        case "open":
            result = OPEN_HOSTPORTSTATUS
        case "filtered":
            result = FILTERED_HOSTPORTSTATUS
        case "closed":
            result = CLOSED_HOSTPORTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HOSTPORTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHostPortStatus(values []HostPortStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HostPortStatus) isMultiValue() bool {
    return false
}
