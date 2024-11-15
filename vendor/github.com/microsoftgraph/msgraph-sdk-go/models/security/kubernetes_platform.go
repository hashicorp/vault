package security
type KubernetesPlatform int

const (
    UNKNOWN_KUBERNETESPLATFORM KubernetesPlatform = iota
    AKS_KUBERNETESPLATFORM
    EKS_KUBERNETESPLATFORM
    GKE_KUBERNETESPLATFORM
    ARC_KUBERNETESPLATFORM
    UNKNOWNFUTUREVALUE_KUBERNETESPLATFORM
)

func (i KubernetesPlatform) String() string {
    return []string{"unknown", "aks", "eks", "gke", "arc", "unknownFutureValue"}[i]
}
func ParseKubernetesPlatform(v string) (any, error) {
    result := UNKNOWN_KUBERNETESPLATFORM
    switch v {
        case "unknown":
            result = UNKNOWN_KUBERNETESPLATFORM
        case "aks":
            result = AKS_KUBERNETESPLATFORM
        case "eks":
            result = EKS_KUBERNETESPLATFORM
        case "gke":
            result = GKE_KUBERNETESPLATFORM
        case "arc":
            result = ARC_KUBERNETESPLATFORM
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_KUBERNETESPLATFORM
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeKubernetesPlatform(values []KubernetesPlatform) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i KubernetesPlatform) isMultiValue() bool {
    return false
}
