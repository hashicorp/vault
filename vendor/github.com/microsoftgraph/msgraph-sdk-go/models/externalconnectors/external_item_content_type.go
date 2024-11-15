package externalconnectors
type ExternalItemContentType int

const (
    TEXT_EXTERNALITEMCONTENTTYPE ExternalItemContentType = iota
    HTML_EXTERNALITEMCONTENTTYPE
    UNKNOWNFUTUREVALUE_EXTERNALITEMCONTENTTYPE
)

func (i ExternalItemContentType) String() string {
    return []string{"text", "html", "unknownFutureValue"}[i]
}
func ParseExternalItemContentType(v string) (any, error) {
    result := TEXT_EXTERNALITEMCONTENTTYPE
    switch v {
        case "text":
            result = TEXT_EXTERNALITEMCONTENTTYPE
        case "html":
            result = HTML_EXTERNALITEMCONTENTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EXTERNALITEMCONTENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeExternalItemContentType(values []ExternalItemContentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ExternalItemContentType) isMultiValue() bool {
    return false
}
