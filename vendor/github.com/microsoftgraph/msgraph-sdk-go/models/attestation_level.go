package models
type AttestationLevel int

const (
    ATTESTED_ATTESTATIONLEVEL AttestationLevel = iota
    NOTATTESTED_ATTESTATIONLEVEL
    UNKNOWNFUTUREVALUE_ATTESTATIONLEVEL
)

func (i AttestationLevel) String() string {
    return []string{"attested", "notAttested", "unknownFutureValue"}[i]
}
func ParseAttestationLevel(v string) (any, error) {
    result := ATTESTED_ATTESTATIONLEVEL
    switch v {
        case "attested":
            result = ATTESTED_ATTESTATIONLEVEL
        case "notAttested":
            result = NOTATTESTED_ATTESTATIONLEVEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ATTESTATIONLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAttestationLevel(values []AttestationLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AttestationLevel) isMultiValue() bool {
    return false
}
