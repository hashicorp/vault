package models
type SelectionLikelihoodInfo int

const (
    NOTSPECIFIED_SELECTIONLIKELIHOODINFO SelectionLikelihoodInfo = iota
    HIGH_SELECTIONLIKELIHOODINFO
)

func (i SelectionLikelihoodInfo) String() string {
    return []string{"notSpecified", "high"}[i]
}
func ParseSelectionLikelihoodInfo(v string) (any, error) {
    result := NOTSPECIFIED_SELECTIONLIKELIHOODINFO
    switch v {
        case "notSpecified":
            result = NOTSPECIFIED_SELECTIONLIKELIHOODINFO
        case "high":
            result = HIGH_SELECTIONLIKELIHOODINFO
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSelectionLikelihoodInfo(values []SelectionLikelihoodInfo) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SelectionLikelihoodInfo) isMultiValue() bool {
    return false
}
