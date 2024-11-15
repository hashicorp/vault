package models
type TitleAreaLayoutType int

const (
    IMAGEANDTITLE_TITLEAREALAYOUTTYPE TitleAreaLayoutType = iota
    PLAIN_TITLEAREALAYOUTTYPE
    COLORBLOCK_TITLEAREALAYOUTTYPE
    OVERLAP_TITLEAREALAYOUTTYPE
    UNKNOWNFUTUREVALUE_TITLEAREALAYOUTTYPE
)

func (i TitleAreaLayoutType) String() string {
    return []string{"imageAndTitle", "plain", "colorBlock", "overlap", "unknownFutureValue"}[i]
}
func ParseTitleAreaLayoutType(v string) (any, error) {
    result := IMAGEANDTITLE_TITLEAREALAYOUTTYPE
    switch v {
        case "imageAndTitle":
            result = IMAGEANDTITLE_TITLEAREALAYOUTTYPE
        case "plain":
            result = PLAIN_TITLEAREALAYOUTTYPE
        case "colorBlock":
            result = COLORBLOCK_TITLEAREALAYOUTTYPE
        case "overlap":
            result = OVERLAP_TITLEAREALAYOUTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TITLEAREALAYOUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTitleAreaLayoutType(values []TitleAreaLayoutType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TitleAreaLayoutType) isMultiValue() bool {
    return false
}
