package models
type BodyType int

const (
    TEXT_BODYTYPE BodyType = iota
    HTML_BODYTYPE
)

func (i BodyType) String() string {
    return []string{"text", "html"}[i]
}
func ParseBodyType(v string) (any, error) {
    result := TEXT_BODYTYPE
    switch v {
        case "text":
            result = TEXT_BODYTYPE
        case "html":
            result = HTML_BODYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBodyType(values []BodyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BodyType) isMultiValue() bool {
    return false
}
