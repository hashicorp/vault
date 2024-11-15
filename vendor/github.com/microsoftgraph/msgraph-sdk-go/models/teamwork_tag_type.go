package models
type TeamworkTagType int

const (
    STANDARD_TEAMWORKTAGTYPE TeamworkTagType = iota
    UNKNOWNFUTUREVALUE_TEAMWORKTAGTYPE
)

func (i TeamworkTagType) String() string {
    return []string{"standard", "unknownFutureValue"}[i]
}
func ParseTeamworkTagType(v string) (any, error) {
    result := STANDARD_TEAMWORKTAGTYPE
    switch v {
        case "standard":
            result = STANDARD_TEAMWORKTAGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TEAMWORKTAGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamworkTagType(values []TeamworkTagType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamworkTagType) isMultiValue() bool {
    return false
}
