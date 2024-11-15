package callrecords
type NetworkTransportProtocol int

const (
    UNKNOWN_NETWORKTRANSPORTPROTOCOL NetworkTransportProtocol = iota
    UDP_NETWORKTRANSPORTPROTOCOL
    TCP_NETWORKTRANSPORTPROTOCOL
    UNKNOWNFUTUREVALUE_NETWORKTRANSPORTPROTOCOL
)

func (i NetworkTransportProtocol) String() string {
    return []string{"unknown", "udp", "tcp", "unknownFutureValue"}[i]
}
func ParseNetworkTransportProtocol(v string) (any, error) {
    result := UNKNOWN_NETWORKTRANSPORTPROTOCOL
    switch v {
        case "unknown":
            result = UNKNOWN_NETWORKTRANSPORTPROTOCOL
        case "udp":
            result = UDP_NETWORKTRANSPORTPROTOCOL
        case "tcp":
            result = TCP_NETWORKTRANSPORTPROTOCOL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_NETWORKTRANSPORTPROTOCOL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeNetworkTransportProtocol(values []NetworkTransportProtocol) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i NetworkTransportProtocol) isMultiValue() bool {
    return false
}
