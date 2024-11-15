package security
type PurgeType int

const (
    RECOVERABLE_PURGETYPE PurgeType = iota
    PERMANENTLYDELETED_PURGETYPE
    UNKNOWNFUTUREVALUE_PURGETYPE
)

func (i PurgeType) String() string {
    return []string{"recoverable", "permanentlyDeleted", "unknownFutureValue"}[i]
}
func ParsePurgeType(v string) (any, error) {
    result := RECOVERABLE_PURGETYPE
    switch v {
        case "recoverable":
            result = RECOVERABLE_PURGETYPE
        case "permanentlyDeleted":
            result = PERMANENTLYDELETED_PURGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PURGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePurgeType(values []PurgeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PurgeType) isMultiValue() bool {
    return false
}
