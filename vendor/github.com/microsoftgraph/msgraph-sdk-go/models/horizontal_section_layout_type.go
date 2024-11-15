package models
type HorizontalSectionLayoutType int

const (
    NONE_HORIZONTALSECTIONLAYOUTTYPE HorizontalSectionLayoutType = iota
    ONECOLUMN_HORIZONTALSECTIONLAYOUTTYPE
    TWOCOLUMNS_HORIZONTALSECTIONLAYOUTTYPE
    THREECOLUMNS_HORIZONTALSECTIONLAYOUTTYPE
    ONETHIRDLEFTCOLUMN_HORIZONTALSECTIONLAYOUTTYPE
    ONETHIRDRIGHTCOLUMN_HORIZONTALSECTIONLAYOUTTYPE
    FULLWIDTH_HORIZONTALSECTIONLAYOUTTYPE
    UNKNOWNFUTUREVALUE_HORIZONTALSECTIONLAYOUTTYPE
)

func (i HorizontalSectionLayoutType) String() string {
    return []string{"none", "oneColumn", "twoColumns", "threeColumns", "oneThirdLeftColumn", "oneThirdRightColumn", "fullWidth", "unknownFutureValue"}[i]
}
func ParseHorizontalSectionLayoutType(v string) (any, error) {
    result := NONE_HORIZONTALSECTIONLAYOUTTYPE
    switch v {
        case "none":
            result = NONE_HORIZONTALSECTIONLAYOUTTYPE
        case "oneColumn":
            result = ONECOLUMN_HORIZONTALSECTIONLAYOUTTYPE
        case "twoColumns":
            result = TWOCOLUMNS_HORIZONTALSECTIONLAYOUTTYPE
        case "threeColumns":
            result = THREECOLUMNS_HORIZONTALSECTIONLAYOUTTYPE
        case "oneThirdLeftColumn":
            result = ONETHIRDLEFTCOLUMN_HORIZONTALSECTIONLAYOUTTYPE
        case "oneThirdRightColumn":
            result = ONETHIRDRIGHTCOLUMN_HORIZONTALSECTIONLAYOUTTYPE
        case "fullWidth":
            result = FULLWIDTH_HORIZONTALSECTIONLAYOUTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HORIZONTALSECTIONLAYOUTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHorizontalSectionLayoutType(values []HorizontalSectionLayoutType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HorizontalSectionLayoutType) isMultiValue() bool {
    return false
}
