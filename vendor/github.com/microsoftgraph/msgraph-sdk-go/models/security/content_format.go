package security
type ContentFormat int

const (
    TEXT_CONTENTFORMAT ContentFormat = iota
    HTML_CONTENTFORMAT
    MARKDOWN_CONTENTFORMAT
    UNKNOWNFUTUREVALUE_CONTENTFORMAT
)

func (i ContentFormat) String() string {
    return []string{"text", "html", "markdown", "unknownFutureValue"}[i]
}
func ParseContentFormat(v string) (any, error) {
    result := TEXT_CONTENTFORMAT
    switch v {
        case "text":
            result = TEXT_CONTENTFORMAT
        case "html":
            result = HTML_CONTENTFORMAT
        case "markdown":
            result = MARKDOWN_CONTENTFORMAT
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONTENTFORMAT
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeContentFormat(values []ContentFormat) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ContentFormat) isMultiValue() bool {
    return false
}
