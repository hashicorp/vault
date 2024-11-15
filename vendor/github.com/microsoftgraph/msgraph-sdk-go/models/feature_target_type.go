package models
type FeatureTargetType int

const (
    GROUP_FEATURETARGETTYPE FeatureTargetType = iota
    ADMINISTRATIVEUNIT_FEATURETARGETTYPE
    ROLE_FEATURETARGETTYPE
    UNKNOWNFUTUREVALUE_FEATURETARGETTYPE
)

func (i FeatureTargetType) String() string {
    return []string{"group", "administrativeUnit", "role", "unknownFutureValue"}[i]
}
func ParseFeatureTargetType(v string) (any, error) {
    result := GROUP_FEATURETARGETTYPE
    switch v {
        case "group":
            result = GROUP_FEATURETARGETTYPE
        case "administrativeUnit":
            result = ADMINISTRATIVEUNIT_FEATURETARGETTYPE
        case "role":
            result = ROLE_FEATURETARGETTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FEATURETARGETTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFeatureTargetType(values []FeatureTargetType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FeatureTargetType) isMultiValue() bool {
    return false
}
