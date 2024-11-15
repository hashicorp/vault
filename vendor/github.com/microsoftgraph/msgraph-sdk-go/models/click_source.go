package models
type ClickSource int

const (
    UNKNOWN_CLICKSOURCE ClickSource = iota
    QRCODE_CLICKSOURCE
    PHISHINGURL_CLICKSOURCE
    UNKNOWNFUTUREVALUE_CLICKSOURCE
)

func (i ClickSource) String() string {
    return []string{"unknown", "qrCode", "phishingUrl", "unknownFutureValue"}[i]
}
func ParseClickSource(v string) (any, error) {
    result := UNKNOWN_CLICKSOURCE
    switch v {
        case "unknown":
            result = UNKNOWN_CLICKSOURCE
        case "qrCode":
            result = QRCODE_CLICKSOURCE
        case "phishingUrl":
            result = PHISHINGURL_CLICKSOURCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLICKSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeClickSource(values []ClickSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ClickSource) isMultiValue() bool {
    return false
}
