package security
type EvidenceRemediationStatus int

const (
    NONE_EVIDENCEREMEDIATIONSTATUS EvidenceRemediationStatus = iota
    REMEDIATED_EVIDENCEREMEDIATIONSTATUS
    PREVENTED_EVIDENCEREMEDIATIONSTATUS
    BLOCKED_EVIDENCEREMEDIATIONSTATUS
    NOTFOUND_EVIDENCEREMEDIATIONSTATUS
    UNKNOWNFUTUREVALUE_EVIDENCEREMEDIATIONSTATUS
)

func (i EvidenceRemediationStatus) String() string {
    return []string{"none", "remediated", "prevented", "blocked", "notFound", "unknownFutureValue"}[i]
}
func ParseEvidenceRemediationStatus(v string) (any, error) {
    result := NONE_EVIDENCEREMEDIATIONSTATUS
    switch v {
        case "none":
            result = NONE_EVIDENCEREMEDIATIONSTATUS
        case "remediated":
            result = REMEDIATED_EVIDENCEREMEDIATIONSTATUS
        case "prevented":
            result = PREVENTED_EVIDENCEREMEDIATIONSTATUS
        case "blocked":
            result = BLOCKED_EVIDENCEREMEDIATIONSTATUS
        case "notFound":
            result = NOTFOUND_EVIDENCEREMEDIATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EVIDENCEREMEDIATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEvidenceRemediationStatus(values []EvidenceRemediationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EvidenceRemediationStatus) isMultiValue() bool {
    return false
}
