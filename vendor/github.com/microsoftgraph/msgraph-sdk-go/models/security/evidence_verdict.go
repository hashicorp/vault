package security
type EvidenceVerdict int

const (
    UNKNOWN_EVIDENCEVERDICT EvidenceVerdict = iota
    SUSPICIOUS_EVIDENCEVERDICT
    MALICIOUS_EVIDENCEVERDICT
    NOTHREATSFOUND_EVIDENCEVERDICT
    UNKNOWNFUTUREVALUE_EVIDENCEVERDICT
)

func (i EvidenceVerdict) String() string {
    return []string{"unknown", "suspicious", "malicious", "noThreatsFound", "unknownFutureValue"}[i]
}
func ParseEvidenceVerdict(v string) (any, error) {
    result := UNKNOWN_EVIDENCEVERDICT
    switch v {
        case "unknown":
            result = UNKNOWN_EVIDENCEVERDICT
        case "suspicious":
            result = SUSPICIOUS_EVIDENCEVERDICT
        case "malicious":
            result = MALICIOUS_EVIDENCEVERDICT
        case "noThreatsFound":
            result = NOTHREATSFOUND_EVIDENCEVERDICT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EVIDENCEVERDICT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEvidenceVerdict(values []EvidenceVerdict) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EvidenceVerdict) isMultiValue() bool {
    return false
}
