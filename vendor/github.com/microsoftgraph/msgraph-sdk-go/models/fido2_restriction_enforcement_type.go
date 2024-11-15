package models
type Fido2RestrictionEnforcementType int

const (
    ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE Fido2RestrictionEnforcementType = iota
    BLOCK_FIDO2RESTRICTIONENFORCEMENTTYPE
    UNKNOWNFUTUREVALUE_FIDO2RESTRICTIONENFORCEMENTTYPE
)

func (i Fido2RestrictionEnforcementType) String() string {
    return []string{"allow", "block", "unknownFutureValue"}[i]
}
func ParseFido2RestrictionEnforcementType(v string) (any, error) {
    result := ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE
    switch v {
        case "allow":
            result = ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE
        case "block":
            result = BLOCK_FIDO2RESTRICTIONENFORCEMENTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FIDO2RESTRICTIONENFORCEMENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFido2RestrictionEnforcementType(values []Fido2RestrictionEnforcementType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Fido2RestrictionEnforcementType) isMultiValue() bool {
    return false
}
