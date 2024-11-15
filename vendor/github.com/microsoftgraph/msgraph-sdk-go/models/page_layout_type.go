package models
type PageLayoutType int

const (
    MICROSOFTRESERVED_PAGELAYOUTTYPE PageLayoutType = iota
    ARTICLE_PAGELAYOUTTYPE
    HOME_PAGELAYOUTTYPE
    UNKNOWNFUTUREVALUE_PAGELAYOUTTYPE
)

func (i PageLayoutType) String() string {
    return []string{"microsoftReserved", "article", "home", "unknownFutureValue"}[i]
}
func ParsePageLayoutType(v string) (any, error) {
    result := MICROSOFTRESERVED_PAGELAYOUTTYPE
    switch v {
        case "microsoftReserved":
            result = MICROSOFTRESERVED_PAGELAYOUTTYPE
        case "article":
            result = ARTICLE_PAGELAYOUTTYPE
        case "home":
            result = HOME_PAGELAYOUTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PAGELAYOUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePageLayoutType(values []PageLayoutType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PageLayoutType) isMultiValue() bool {
    return false
}
