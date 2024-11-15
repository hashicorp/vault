package models
type RoutingType int

const (
    FORWARDED_ROUTINGTYPE RoutingType = iota
    LOOKUP_ROUTINGTYPE
    SELFFORK_ROUTINGTYPE
    UNKNOWNFUTUREVALUE_ROUTINGTYPE
)

func (i RoutingType) String() string {
    return []string{"forwarded", "lookup", "selfFork", "unknownFutureValue"}[i]
}
func ParseRoutingType(v string) (any, error) {
    result := FORWARDED_ROUTINGTYPE
    switch v {
        case "forwarded":
            result = FORWARDED_ROUTINGTYPE
        case "lookup":
            result = LOOKUP_ROUTINGTYPE
        case "selfFork":
            result = SELFFORK_ROUTINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ROUTINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRoutingType(values []RoutingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RoutingType) isMultiValue() bool {
    return false
}
