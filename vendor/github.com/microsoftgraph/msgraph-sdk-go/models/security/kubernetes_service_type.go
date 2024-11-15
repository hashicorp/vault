package security
type KubernetesServiceType int

const (
    UNKNOWN_KUBERNETESSERVICETYPE KubernetesServiceType = iota
    CLUSTERIP_KUBERNETESSERVICETYPE
    EXTERNALNAME_KUBERNETESSERVICETYPE
    NODEPORT_KUBERNETESSERVICETYPE
    LOADBALANCER_KUBERNETESSERVICETYPE
    UNKNOWNFUTUREVALUE_KUBERNETESSERVICETYPE
)

func (i KubernetesServiceType) String() string {
    return []string{"unknown", "clusterIP", "externalName", "nodePort", "loadBalancer", "unknownFutureValue"}[i]
}
func ParseKubernetesServiceType(v string) (any, error) {
    result := UNKNOWN_KUBERNETESSERVICETYPE
    switch v {
        case "unknown":
            result = UNKNOWN_KUBERNETESSERVICETYPE
        case "clusterIP":
            result = CLUSTERIP_KUBERNETESSERVICETYPE
        case "externalName":
            result = EXTERNALNAME_KUBERNETESSERVICETYPE
        case "nodePort":
            result = NODEPORT_KUBERNETESSERVICETYPE
        case "loadBalancer":
            result = LOADBALANCER_KUBERNETESSERVICETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_KUBERNETESSERVICETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeKubernetesServiceType(values []KubernetesServiceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i KubernetesServiceType) isMultiValue() bool {
    return false
}
