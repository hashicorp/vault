package models
type ServiceHealthOrigin int

const (
    MICROSOFT_SERVICEHEALTHORIGIN ServiceHealthOrigin = iota
    THIRDPARTY_SERVICEHEALTHORIGIN
    CUSTOMER_SERVICEHEALTHORIGIN
    UNKNOWNFUTUREVALUE_SERVICEHEALTHORIGIN
)

func (i ServiceHealthOrigin) String() string {
    return []string{"microsoft", "thirdParty", "customer", "unknownFutureValue"}[i]
}
func ParseServiceHealthOrigin(v string) (any, error) {
    result := MICROSOFT_SERVICEHEALTHORIGIN
    switch v {
        case "microsoft":
            result = MICROSOFT_SERVICEHEALTHORIGIN
        case "thirdParty":
            result = THIRDPARTY_SERVICEHEALTHORIGIN
        case "customer":
            result = CUSTOMER_SERVICEHEALTHORIGIN
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICEHEALTHORIGIN
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServiceHealthOrigin(values []ServiceHealthOrigin) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServiceHealthOrigin) isMultiValue() bool {
    return false
}
