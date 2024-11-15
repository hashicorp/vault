package models
type AccountTargetContentType int

const (
    UNKNOWN_ACCOUNTTARGETCONTENTTYPE AccountTargetContentType = iota
    INCLUDEALL_ACCOUNTTARGETCONTENTTYPE
    ADDRESSBOOK_ACCOUNTTARGETCONTENTTYPE
    UNKNOWNFUTUREVALUE_ACCOUNTTARGETCONTENTTYPE
)

func (i AccountTargetContentType) String() string {
    return []string{"unknown", "includeAll", "addressBook", "unknownFutureValue"}[i]
}
func ParseAccountTargetContentType(v string) (any, error) {
    result := UNKNOWN_ACCOUNTTARGETCONTENTTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_ACCOUNTTARGETCONTENTTYPE
        case "includeAll":
            result = INCLUDEALL_ACCOUNTTARGETCONTENTTYPE
        case "addressBook":
            result = ADDRESSBOOK_ACCOUNTTARGETCONTENTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ACCOUNTTARGETCONTENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAccountTargetContentType(values []AccountTargetContentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AccountTargetContentType) isMultiValue() bool {
    return false
}
