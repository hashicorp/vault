package models
type CategoryColor int

const (
    NONE_CATEGORYCOLOR CategoryColor = iota
    PRESET0_CATEGORYCOLOR
    PRESET1_CATEGORYCOLOR
    PRESET2_CATEGORYCOLOR
    PRESET3_CATEGORYCOLOR
    PRESET4_CATEGORYCOLOR
    PRESET5_CATEGORYCOLOR
    PRESET6_CATEGORYCOLOR
    PRESET7_CATEGORYCOLOR
    PRESET8_CATEGORYCOLOR
    PRESET9_CATEGORYCOLOR
    PRESET10_CATEGORYCOLOR
    PRESET11_CATEGORYCOLOR
    PRESET12_CATEGORYCOLOR
    PRESET13_CATEGORYCOLOR
    PRESET14_CATEGORYCOLOR
    PRESET15_CATEGORYCOLOR
    PRESET16_CATEGORYCOLOR
    PRESET17_CATEGORYCOLOR
    PRESET18_CATEGORYCOLOR
    PRESET19_CATEGORYCOLOR
    PRESET20_CATEGORYCOLOR
    PRESET21_CATEGORYCOLOR
    PRESET22_CATEGORYCOLOR
    PRESET23_CATEGORYCOLOR
    PRESET24_CATEGORYCOLOR
)

func (i CategoryColor) String() string {
    return []string{"none", "preset0", "preset1", "preset2", "preset3", "preset4", "preset5", "preset6", "preset7", "preset8", "preset9", "preset10", "preset11", "preset12", "preset13", "preset14", "preset15", "preset16", "preset17", "preset18", "preset19", "preset20", "preset21", "preset22", "preset23", "preset24"}[i]
}
func ParseCategoryColor(v string) (any, error) {
    result := NONE_CATEGORYCOLOR
    switch v {
        case "none":
            result = NONE_CATEGORYCOLOR
        case "preset0":
            result = PRESET0_CATEGORYCOLOR
        case "preset1":
            result = PRESET1_CATEGORYCOLOR
        case "preset2":
            result = PRESET2_CATEGORYCOLOR
        case "preset3":
            result = PRESET3_CATEGORYCOLOR
        case "preset4":
            result = PRESET4_CATEGORYCOLOR
        case "preset5":
            result = PRESET5_CATEGORYCOLOR
        case "preset6":
            result = PRESET6_CATEGORYCOLOR
        case "preset7":
            result = PRESET7_CATEGORYCOLOR
        case "preset8":
            result = PRESET8_CATEGORYCOLOR
        case "preset9":
            result = PRESET9_CATEGORYCOLOR
        case "preset10":
            result = PRESET10_CATEGORYCOLOR
        case "preset11":
            result = PRESET11_CATEGORYCOLOR
        case "preset12":
            result = PRESET12_CATEGORYCOLOR
        case "preset13":
            result = PRESET13_CATEGORYCOLOR
        case "preset14":
            result = PRESET14_CATEGORYCOLOR
        case "preset15":
            result = PRESET15_CATEGORYCOLOR
        case "preset16":
            result = PRESET16_CATEGORYCOLOR
        case "preset17":
            result = PRESET17_CATEGORYCOLOR
        case "preset18":
            result = PRESET18_CATEGORYCOLOR
        case "preset19":
            result = PRESET19_CATEGORYCOLOR
        case "preset20":
            result = PRESET20_CATEGORYCOLOR
        case "preset21":
            result = PRESET21_CATEGORYCOLOR
        case "preset22":
            result = PRESET22_CATEGORYCOLOR
        case "preset23":
            result = PRESET23_CATEGORYCOLOR
        case "preset24":
            result = PRESET24_CATEGORYCOLOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCategoryColor(values []CategoryColor) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CategoryColor) isMultiValue() bool {
    return false
}
