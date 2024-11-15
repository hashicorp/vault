package models
type QuarantineReason int

const (
    ENCOUNTEREDBASEESCROWTHRESHOLD_QUARANTINEREASON QuarantineReason = iota
    ENCOUNTEREDTOTALESCROWTHRESHOLD_QUARANTINEREASON
    ENCOUNTEREDESCROWPROPORTIONTHRESHOLD_QUARANTINEREASON
    ENCOUNTEREDQUARANTINEEXCEPTION_QUARANTINEREASON
    UNKNOWN_QUARANTINEREASON
    QUARANTINEDONDEMAND_QUARANTINEREASON
    TOOMANYDELETES_QUARANTINEREASON
    INGESTIONINTERRUPTED_QUARANTINEREASON
)

func (i QuarantineReason) String() string {
    return []string{"EncounteredBaseEscrowThreshold", "EncounteredTotalEscrowThreshold", "EncounteredEscrowProportionThreshold", "EncounteredQuarantineException", "Unknown", "QuarantinedOnDemand", "TooManyDeletes", "IngestionInterrupted"}[i]
}
func ParseQuarantineReason(v string) (any, error) {
    result := ENCOUNTEREDBASEESCROWTHRESHOLD_QUARANTINEREASON
    switch v {
        case "EncounteredBaseEscrowThreshold":
            result = ENCOUNTEREDBASEESCROWTHRESHOLD_QUARANTINEREASON
        case "EncounteredTotalEscrowThreshold":
            result = ENCOUNTEREDTOTALESCROWTHRESHOLD_QUARANTINEREASON
        case "EncounteredEscrowProportionThreshold":
            result = ENCOUNTEREDESCROWPROPORTIONTHRESHOLD_QUARANTINEREASON
        case "EncounteredQuarantineException":
            result = ENCOUNTEREDQUARANTINEEXCEPTION_QUARANTINEREASON
        case "Unknown":
            result = UNKNOWN_QUARANTINEREASON
        case "QuarantinedOnDemand":
            result = QUARANTINEDONDEMAND_QUARANTINEREASON
        case "TooManyDeletes":
            result = TOOMANYDELETES_QUARANTINEREASON
        case "IngestionInterrupted":
            result = INGESTIONINTERRUPTED_QUARANTINEREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeQuarantineReason(values []QuarantineReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i QuarantineReason) isMultiValue() bool {
    return false
}
