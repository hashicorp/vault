package models
type BucketAggregationSortProperty int

const (
    COUNT_BUCKETAGGREGATIONSORTPROPERTY BucketAggregationSortProperty = iota
    KEYASSTRING_BUCKETAGGREGATIONSORTPROPERTY
    KEYASNUMBER_BUCKETAGGREGATIONSORTPROPERTY
    UNKNOWNFUTUREVALUE_BUCKETAGGREGATIONSORTPROPERTY
)

func (i BucketAggregationSortProperty) String() string {
    return []string{"count", "keyAsString", "keyAsNumber", "unknownFutureValue"}[i]
}
func ParseBucketAggregationSortProperty(v string) (any, error) {
    result := COUNT_BUCKETAGGREGATIONSORTPROPERTY
    switch v {
        case "count":
            result = COUNT_BUCKETAGGREGATIONSORTPROPERTY
        case "keyAsString":
            result = KEYASSTRING_BUCKETAGGREGATIONSORTPROPERTY
        case "keyAsNumber":
            result = KEYASNUMBER_BUCKETAGGREGATIONSORTPROPERTY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BUCKETAGGREGATIONSORTPROPERTY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBucketAggregationSortProperty(values []BucketAggregationSortProperty) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BucketAggregationSortProperty) isMultiValue() bool {
    return false
}
