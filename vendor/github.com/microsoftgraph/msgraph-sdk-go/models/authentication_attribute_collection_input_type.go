package models
type AuthenticationAttributeCollectionInputType int

const (
    TEXT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE AuthenticationAttributeCollectionInputType = iota
    RADIOSINGLESELECT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
    CHECKBOXMULTISELECT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
    BOOLEAN_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
)

func (i AuthenticationAttributeCollectionInputType) String() string {
    return []string{"text", "radioSingleSelect", "checkboxMultiSelect", "boolean", "unknownFutureValue"}[i]
}
func ParseAuthenticationAttributeCollectionInputType(v string) (any, error) {
    result := TEXT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
    switch v {
        case "text":
            result = TEXT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
        case "radioSingleSelect":
            result = RADIOSINGLESELECT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
        case "checkboxMultiSelect":
            result = CHECKBOXMULTISELECT_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
        case "boolean":
            result = BOOLEAN_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONATTRIBUTECOLLECTIONINPUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationAttributeCollectionInputType(values []AuthenticationAttributeCollectionInputType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationAttributeCollectionInputType) isMultiValue() bool {
    return false
}
