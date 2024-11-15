package models
type EndpointType int

const (
    DEFAULT_ENDPOINTTYPE EndpointType = iota
    VOICEMAIL_ENDPOINTTYPE
    SKYPEFORBUSINESS_ENDPOINTTYPE
    SKYPEFORBUSINESSVOIPPHONE_ENDPOINTTYPE
    UNKNOWNFUTUREVALUE_ENDPOINTTYPE
)

func (i EndpointType) String() string {
    return []string{"default", "voicemail", "skypeForBusiness", "skypeForBusinessVoipPhone", "unknownFutureValue"}[i]
}
func ParseEndpointType(v string) (any, error) {
    result := DEFAULT_ENDPOINTTYPE
    switch v {
        case "default":
            result = DEFAULT_ENDPOINTTYPE
        case "voicemail":
            result = VOICEMAIL_ENDPOINTTYPE
        case "skypeForBusiness":
            result = SKYPEFORBUSINESS_ENDPOINTTYPE
        case "skypeForBusinessVoipPhone":
            result = SKYPEFORBUSINESSVOIPPHONE_ENDPOINTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ENDPOINTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEndpointType(values []EndpointType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EndpointType) isMultiValue() bool {
    return false
}
