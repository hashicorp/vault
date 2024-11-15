package billing
type BillingPeriod int

const (
    CURRENT_BILLINGPERIOD BillingPeriod = iota
    LAST_BILLINGPERIOD
    UNKNOWNFUTUREVALUE_BILLINGPERIOD
)

func (i BillingPeriod) String() string {
    return []string{"current", "last", "unknownFutureValue"}[i]
}
func ParseBillingPeriod(v string) (any, error) {
    result := CURRENT_BILLINGPERIOD
    switch v {
        case "current":
            result = CURRENT_BILLINGPERIOD
        case "last":
            result = LAST_BILLINGPERIOD
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BILLINGPERIOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBillingPeriod(values []BillingPeriod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BillingPeriod) isMultiValue() bool {
    return false
}
