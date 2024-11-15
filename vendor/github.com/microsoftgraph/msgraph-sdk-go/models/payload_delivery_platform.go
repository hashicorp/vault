package models
type PayloadDeliveryPlatform int

const (
    UNKNOWN_PAYLOADDELIVERYPLATFORM PayloadDeliveryPlatform = iota
    SMS_PAYLOADDELIVERYPLATFORM
    EMAIL_PAYLOADDELIVERYPLATFORM
    TEAMS_PAYLOADDELIVERYPLATFORM
    UNKNOWNFUTUREVALUE_PAYLOADDELIVERYPLATFORM
)

func (i PayloadDeliveryPlatform) String() string {
    return []string{"unknown", "sms", "email", "teams", "unknownFutureValue"}[i]
}
func ParsePayloadDeliveryPlatform(v string) (any, error) {
    result := UNKNOWN_PAYLOADDELIVERYPLATFORM
    switch v {
        case "unknown":
            result = UNKNOWN_PAYLOADDELIVERYPLATFORM
        case "sms":
            result = SMS_PAYLOADDELIVERYPLATFORM
        case "email":
            result = EMAIL_PAYLOADDELIVERYPLATFORM
        case "teams":
            result = TEAMS_PAYLOADDELIVERYPLATFORM
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PAYLOADDELIVERYPLATFORM
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePayloadDeliveryPlatform(values []PayloadDeliveryPlatform) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PayloadDeliveryPlatform) isMultiValue() bool {
    return false
}
