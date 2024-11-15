package models
type AppsUpdateChannelType int

const (
    CURRENT_APPSUPDATECHANNELTYPE AppsUpdateChannelType = iota
    MONTHLYENTERPRISE_APPSUPDATECHANNELTYPE
    SEMIANNUAL_APPSUPDATECHANNELTYPE
    UNKNOWNFUTUREVALUE_APPSUPDATECHANNELTYPE
)

func (i AppsUpdateChannelType) String() string {
    return []string{"current", "monthlyEnterprise", "semiAnnual", "unknownFutureValue"}[i]
}
func ParseAppsUpdateChannelType(v string) (any, error) {
    result := CURRENT_APPSUPDATECHANNELTYPE
    switch v {
        case "current":
            result = CURRENT_APPSUPDATECHANNELTYPE
        case "monthlyEnterprise":
            result = MONTHLYENTERPRISE_APPSUPDATECHANNELTYPE
        case "semiAnnual":
            result = SEMIANNUAL_APPSUPDATECHANNELTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_APPSUPDATECHANNELTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppsUpdateChannelType(values []AppsUpdateChannelType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppsUpdateChannelType) isMultiValue() bool {
    return false
}
