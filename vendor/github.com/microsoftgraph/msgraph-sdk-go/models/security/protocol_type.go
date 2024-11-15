package security
type ProtocolType int

const (
    TCP_PROTOCOLTYPE ProtocolType = iota
    UDP_PROTOCOLTYPE
    UNKNOWNFUTUREVALUE_PROTOCOLTYPE
)

func (i ProtocolType) String() string {
    return []string{"tcp", "udp", "unknownFutureValue"}[i]
}
func ParseProtocolType(v string) (any, error) {
    result := TCP_PROTOCOLTYPE
    switch v {
        case "tcp":
            result = TCP_PROTOCOLTYPE
        case "udp":
            result = UDP_PROTOCOLTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROTOCOLTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProtocolType(values []ProtocolType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProtocolType) isMultiValue() bool {
    return false
}
