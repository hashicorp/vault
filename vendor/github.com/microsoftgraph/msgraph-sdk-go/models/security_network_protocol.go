package models
type SecurityNetworkProtocol int

const (
    UNKNOWN_SECURITYNETWORKPROTOCOL SecurityNetworkProtocol = iota
    IP_SECURITYNETWORKPROTOCOL
    ICMP_SECURITYNETWORKPROTOCOL
    IGMP_SECURITYNETWORKPROTOCOL
    GGP_SECURITYNETWORKPROTOCOL
    IPV4_SECURITYNETWORKPROTOCOL
    TCP_SECURITYNETWORKPROTOCOL
    PUP_SECURITYNETWORKPROTOCOL
    UDP_SECURITYNETWORKPROTOCOL
    IDP_SECURITYNETWORKPROTOCOL
    IPV6_SECURITYNETWORKPROTOCOL
    IPV6ROUTINGHEADER_SECURITYNETWORKPROTOCOL
    IPV6FRAGMENTHEADER_SECURITYNETWORKPROTOCOL
    IPSECENCAPSULATINGSECURITYPAYLOAD_SECURITYNETWORKPROTOCOL
    IPSECAUTHENTICATIONHEADER_SECURITYNETWORKPROTOCOL
    ICMPV6_SECURITYNETWORKPROTOCOL
    IPV6NONEXTHEADER_SECURITYNETWORKPROTOCOL
    IPV6DESTINATIONOPTIONS_SECURITYNETWORKPROTOCOL
    ND_SECURITYNETWORKPROTOCOL
    RAW_SECURITYNETWORKPROTOCOL
    IPX_SECURITYNETWORKPROTOCOL
    SPX_SECURITYNETWORKPROTOCOL
    SPXII_SECURITYNETWORKPROTOCOL
    UNKNOWNFUTUREVALUE_SECURITYNETWORKPROTOCOL
)

func (i SecurityNetworkProtocol) String() string {
    return []string{"unknown", "ip", "icmp", "igmp", "ggp", "ipv4", "tcp", "pup", "udp", "idp", "ipv6", "ipv6RoutingHeader", "ipv6FragmentHeader", "ipSecEncapsulatingSecurityPayload", "ipSecAuthenticationHeader", "icmpV6", "ipv6NoNextHeader", "ipv6DestinationOptions", "nd", "raw", "ipx", "spx", "spxII", "unknownFutureValue"}[i]
}
func ParseSecurityNetworkProtocol(v string) (any, error) {
    result := UNKNOWN_SECURITYNETWORKPROTOCOL
    switch v {
        case "unknown":
            result = UNKNOWN_SECURITYNETWORKPROTOCOL
        case "ip":
            result = IP_SECURITYNETWORKPROTOCOL
        case "icmp":
            result = ICMP_SECURITYNETWORKPROTOCOL
        case "igmp":
            result = IGMP_SECURITYNETWORKPROTOCOL
        case "ggp":
            result = GGP_SECURITYNETWORKPROTOCOL
        case "ipv4":
            result = IPV4_SECURITYNETWORKPROTOCOL
        case "tcp":
            result = TCP_SECURITYNETWORKPROTOCOL
        case "pup":
            result = PUP_SECURITYNETWORKPROTOCOL
        case "udp":
            result = UDP_SECURITYNETWORKPROTOCOL
        case "idp":
            result = IDP_SECURITYNETWORKPROTOCOL
        case "ipv6":
            result = IPV6_SECURITYNETWORKPROTOCOL
        case "ipv6RoutingHeader":
            result = IPV6ROUTINGHEADER_SECURITYNETWORKPROTOCOL
        case "ipv6FragmentHeader":
            result = IPV6FRAGMENTHEADER_SECURITYNETWORKPROTOCOL
        case "ipSecEncapsulatingSecurityPayload":
            result = IPSECENCAPSULATINGSECURITYPAYLOAD_SECURITYNETWORKPROTOCOL
        case "ipSecAuthenticationHeader":
            result = IPSECAUTHENTICATIONHEADER_SECURITYNETWORKPROTOCOL
        case "icmpV6":
            result = ICMPV6_SECURITYNETWORKPROTOCOL
        case "ipv6NoNextHeader":
            result = IPV6NONEXTHEADER_SECURITYNETWORKPROTOCOL
        case "ipv6DestinationOptions":
            result = IPV6DESTINATIONOPTIONS_SECURITYNETWORKPROTOCOL
        case "nd":
            result = ND_SECURITYNETWORKPROTOCOL
        case "raw":
            result = RAW_SECURITYNETWORKPROTOCOL
        case "ipx":
            result = IPX_SECURITYNETWORKPROTOCOL
        case "spx":
            result = SPX_SECURITYNETWORKPROTOCOL
        case "spxII":
            result = SPXII_SECURITYNETWORKPROTOCOL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SECURITYNETWORKPROTOCOL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSecurityNetworkProtocol(values []SecurityNetworkProtocol) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SecurityNetworkProtocol) isMultiValue() bool {
    return false
}
