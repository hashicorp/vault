package models
type DestinationType int

const (
    NEW_DESTINATIONTYPE DestinationType = iota
    INPLACE_DESTINATIONTYPE
    UNKNOWNFUTUREVALUE_DESTINATIONTYPE
)

func (i DestinationType) String() string {
    return []string{"new", "inPlace", "unknownFutureValue"}[i]
}
func ParseDestinationType(v string) (any, error) {
    result := NEW_DESTINATIONTYPE
    switch v {
        case "new":
            result = NEW_DESTINATIONTYPE
        case "inPlace":
            result = INPLACE_DESTINATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DESTINATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDestinationType(values []DestinationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DestinationType) isMultiValue() bool {
    return false
}
