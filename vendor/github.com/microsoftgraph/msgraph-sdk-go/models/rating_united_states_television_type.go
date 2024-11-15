package models
// TV content rating labels in United States
type RatingUnitedStatesTelevisionType int

const (
    // Default value, allow all TV shows content
    ALLALLOWED_RATINGUNITEDSTATESTELEVISIONTYPE RatingUnitedStatesTelevisionType = iota
    // Do not allow any TV shows content
    ALLBLOCKED_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-Y, all children
    CHILDRENALL_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-Y7, children age 7 and above
    CHILDRENABOVE7_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-G, suitable for all ages
    GENERAL_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-PG, parental guidance
    PARENTALGUIDANCE_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-14, children age 14 and above
    CHILDRENABOVE14_RATINGUNITEDSTATESTELEVISIONTYPE
    // TV-MA, adults only
    ADULTS_RATINGUNITEDSTATESTELEVISIONTYPE
)

func (i RatingUnitedStatesTelevisionType) String() string {
    return []string{"allAllowed", "allBlocked", "childrenAll", "childrenAbove7", "general", "parentalGuidance", "childrenAbove14", "adults"}[i]
}
func ParseRatingUnitedStatesTelevisionType(v string) (any, error) {
    result := ALLALLOWED_RATINGUNITEDSTATESTELEVISIONTYPE
    switch v {
        case "allAllowed":
            result = ALLALLOWED_RATINGUNITEDSTATESTELEVISIONTYPE
        case "allBlocked":
            result = ALLBLOCKED_RATINGUNITEDSTATESTELEVISIONTYPE
        case "childrenAll":
            result = CHILDRENALL_RATINGUNITEDSTATESTELEVISIONTYPE
        case "childrenAbove7":
            result = CHILDRENABOVE7_RATINGUNITEDSTATESTELEVISIONTYPE
        case "general":
            result = GENERAL_RATINGUNITEDSTATESTELEVISIONTYPE
        case "parentalGuidance":
            result = PARENTALGUIDANCE_RATINGUNITEDSTATESTELEVISIONTYPE
        case "childrenAbove14":
            result = CHILDRENABOVE14_RATINGUNITEDSTATESTELEVISIONTYPE
        case "adults":
            result = ADULTS_RATINGUNITEDSTATESTELEVISIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRatingUnitedStatesTelevisionType(values []RatingUnitedStatesTelevisionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RatingUnitedStatesTelevisionType) isMultiValue() bool {
    return false
}
