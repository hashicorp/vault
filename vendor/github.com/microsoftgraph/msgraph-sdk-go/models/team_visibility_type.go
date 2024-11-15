package models
type TeamVisibilityType int

const (
    PRIVATE_TEAMVISIBILITYTYPE TeamVisibilityType = iota
    PUBLIC_TEAMVISIBILITYTYPE
    HIDDENMEMBERSHIP_TEAMVISIBILITYTYPE
    UNKNOWNFUTUREVALUE_TEAMVISIBILITYTYPE
)

func (i TeamVisibilityType) String() string {
    return []string{"private", "public", "hiddenMembership", "unknownFutureValue"}[i]
}
func ParseTeamVisibilityType(v string) (any, error) {
    result := PRIVATE_TEAMVISIBILITYTYPE
    switch v {
        case "private":
            result = PRIVATE_TEAMVISIBILITYTYPE
        case "public":
            result = PUBLIC_TEAMVISIBILITYTYPE
        case "hiddenMembership":
            result = HIDDENMEMBERSHIP_TEAMVISIBILITYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TEAMVISIBILITYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamVisibilityType(values []TeamVisibilityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamVisibilityType) isMultiValue() bool {
    return false
}
