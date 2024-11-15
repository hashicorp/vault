package security
type EvidenceRole int

const (
    UNKNOWN_EVIDENCEROLE EvidenceRole = iota
    CONTEXTUAL_EVIDENCEROLE
    SCANNED_EVIDENCEROLE
    SOURCE_EVIDENCEROLE
    DESTINATION_EVIDENCEROLE
    CREATED_EVIDENCEROLE
    ADDED_EVIDENCEROLE
    COMPROMISED_EVIDENCEROLE
    EDITED_EVIDENCEROLE
    ATTACKED_EVIDENCEROLE
    ATTACKER_EVIDENCEROLE
    COMMANDANDCONTROL_EVIDENCEROLE
    LOADED_EVIDENCEROLE
    SUSPICIOUS_EVIDENCEROLE
    POLICYVIOLATOR_EVIDENCEROLE
    UNKNOWNFUTUREVALUE_EVIDENCEROLE
)

func (i EvidenceRole) String() string {
    return []string{"unknown", "contextual", "scanned", "source", "destination", "created", "added", "compromised", "edited", "attacked", "attacker", "commandAndControl", "loaded", "suspicious", "policyViolator", "unknownFutureValue"}[i]
}
func ParseEvidenceRole(v string) (any, error) {
    result := UNKNOWN_EVIDENCEROLE
    switch v {
        case "unknown":
            result = UNKNOWN_EVIDENCEROLE
        case "contextual":
            result = CONTEXTUAL_EVIDENCEROLE
        case "scanned":
            result = SCANNED_EVIDENCEROLE
        case "source":
            result = SOURCE_EVIDENCEROLE
        case "destination":
            result = DESTINATION_EVIDENCEROLE
        case "created":
            result = CREATED_EVIDENCEROLE
        case "added":
            result = ADDED_EVIDENCEROLE
        case "compromised":
            result = COMPROMISED_EVIDENCEROLE
        case "edited":
            result = EDITED_EVIDENCEROLE
        case "attacked":
            result = ATTACKED_EVIDENCEROLE
        case "attacker":
            result = ATTACKER_EVIDENCEROLE
        case "commandAndControl":
            result = COMMANDANDCONTROL_EVIDENCEROLE
        case "loaded":
            result = LOADED_EVIDENCEROLE
        case "suspicious":
            result = SUSPICIOUS_EVIDENCEROLE
        case "policyViolator":
            result = POLICYVIOLATOR_EVIDENCEROLE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EVIDENCEROLE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEvidenceRole(values []EvidenceRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EvidenceRole) isMultiValue() bool {
    return false
}
