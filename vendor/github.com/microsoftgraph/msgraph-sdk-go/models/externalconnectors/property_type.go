package externalconnectors
type PropertyType int

const (
    STRING_PROPERTYTYPE PropertyType = iota
    INT64_PROPERTYTYPE
    DOUBLE_PROPERTYTYPE
    DATETIME_PROPERTYTYPE
    BOOLEAN_PROPERTYTYPE
    STRINGCOLLECTION_PROPERTYTYPE
    INT64COLLECTION_PROPERTYTYPE
    DOUBLECOLLECTION_PROPERTYTYPE
    DATETIMECOLLECTION_PROPERTYTYPE
    UNKNOWNFUTUREVALUE_PROPERTYTYPE
)

func (i PropertyType) String() string {
    return []string{"string", "int64", "double", "dateTime", "boolean", "stringCollection", "int64Collection", "doubleCollection", "dateTimeCollection", "unknownFutureValue"}[i]
}
func ParsePropertyType(v string) (any, error) {
    result := STRING_PROPERTYTYPE
    switch v {
        case "string":
            result = STRING_PROPERTYTYPE
        case "int64":
            result = INT64_PROPERTYTYPE
        case "double":
            result = DOUBLE_PROPERTYTYPE
        case "dateTime":
            result = DATETIME_PROPERTYTYPE
        case "boolean":
            result = BOOLEAN_PROPERTYTYPE
        case "stringCollection":
            result = STRINGCOLLECTION_PROPERTYTYPE
        case "int64Collection":
            result = INT64COLLECTION_PROPERTYTYPE
        case "doubleCollection":
            result = DOUBLECOLLECTION_PROPERTYTYPE
        case "dateTimeCollection":
            result = DATETIMECOLLECTION_PROPERTYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PROPERTYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePropertyType(values []PropertyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PropertyType) isMultiValue() bool {
    return false
}
