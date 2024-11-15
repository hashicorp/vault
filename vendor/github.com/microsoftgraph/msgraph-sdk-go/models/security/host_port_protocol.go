package security
type HostPortProtocol int

const (
    TCP_HOSTPORTPROTOCOL HostPortProtocol = iota
    UDP_HOSTPORTPROTOCOL
    UNKNOWNFUTUREVALUE_HOSTPORTPROTOCOL
)

func (i HostPortProtocol) String() string {
    return []string{"tcp", "udp", "unknownFutureValue"}[i]
}
func ParseHostPortProtocol(v string) (any, error) {
    result := TCP_HOSTPORTPROTOCOL
    switch v {
        case "tcp":
            result = TCP_HOSTPORTPROTOCOL
        case "udp":
            result = UDP_HOSTPORTPROTOCOL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HOSTPORTPROTOCOL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHostPortProtocol(values []HostPortProtocol) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HostPortProtocol) isMultiValue() bool {
    return false
}
