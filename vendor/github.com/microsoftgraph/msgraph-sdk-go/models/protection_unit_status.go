package models
type ProtectionUnitStatus int

const (
    PROTECTREQUESTED_PROTECTIONUNITSTATUS ProtectionUnitStatus = iota
    PROTECTED_PROTECTIONUNITSTATUS
    UNPROTECTREQUESTED_PROTECTIONUNITSTATUS
    UNPROTECTED_PROTECTIONUNITSTATUS
    REMOVEREQUESTED_PROTECTIONUNITSTATUS
    UNKNOWNFUTUREVALUE_PROTECTIONUNITSTATUS
)

func (i ProtectionUnitStatus) String() string {
    return []string{"protectRequested", "protected", "unprotectRequested", "unprotected", "removeRequested", "unknownFutureValue"}[i]
}
func ParseProtectionUnitStatus(v string) (any, error) {
    result := PROTECTREQUESTED_PROTECTIONUNITSTATUS
    switch v {
        case "protectRequested":
            result = PROTECTREQUESTED_PROTECTIONUNITSTATUS
        case "protected":
            result = PROTECTED_PROTECTIONUNITSTATUS
        case "unprotectRequested":
            result = UNPROTECTREQUESTED_PROTECTIONUNITSTATUS
        case "unprotected":
            result = UNPROTECTED_PROTECTIONUNITSTATUS
        case "removeRequested":
            result = REMOVEREQUESTED_PROTECTIONUNITSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROTECTIONUNITSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProtectionUnitStatus(values []ProtectionUnitStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProtectionUnitStatus) isMultiValue() bool {
    return false
}
