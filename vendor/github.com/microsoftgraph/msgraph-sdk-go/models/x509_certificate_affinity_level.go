package models
type X509CertificateAffinityLevel int

const (
    LOW_X509CERTIFICATEAFFINITYLEVEL X509CertificateAffinityLevel = iota
    HIGH_X509CERTIFICATEAFFINITYLEVEL
    UNKNOWNFUTUREVALUE_X509CERTIFICATEAFFINITYLEVEL
)

func (i X509CertificateAffinityLevel) String() string {
    return []string{"low", "high", "unknownFutureValue"}[i]
}
func ParseX509CertificateAffinityLevel(v string) (any, error) {
    result := LOW_X509CERTIFICATEAFFINITYLEVEL
    switch v {
        case "low":
            result = LOW_X509CERTIFICATEAFFINITYLEVEL
        case "high":
            result = HIGH_X509CERTIFICATEAFFINITYLEVEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_X509CERTIFICATEAFFINITYLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeX509CertificateAffinityLevel(values []X509CertificateAffinityLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i X509CertificateAffinityLevel) isMultiValue() bool {
    return false
}
