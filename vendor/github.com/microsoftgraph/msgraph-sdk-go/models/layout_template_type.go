package models
type LayoutTemplateType int

const (
    DEFAULT_LAYOUTTEMPLATETYPE LayoutTemplateType = iota
    VERTICALSPLIT_LAYOUTTEMPLATETYPE
    UNKNOWNFUTUREVALUE_LAYOUTTEMPLATETYPE
)

func (i LayoutTemplateType) String() string {
    return []string{"default", "verticalSplit", "unknownFutureValue"}[i]
}
func ParseLayoutTemplateType(v string) (any, error) {
    result := DEFAULT_LAYOUTTEMPLATETYPE
    switch v {
        case "default":
            result = DEFAULT_LAYOUTTEMPLATETYPE
        case "verticalSplit":
            result = VERTICALSPLIT_LAYOUTTEMPLATETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_LAYOUTTEMPLATETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLayoutTemplateType(values []LayoutTemplateType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LayoutTemplateType) isMultiValue() bool {
    return false
}
