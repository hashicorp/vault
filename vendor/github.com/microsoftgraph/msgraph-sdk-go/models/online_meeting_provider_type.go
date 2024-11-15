package models
type OnlineMeetingProviderType int

const (
    UNKNOWN_ONLINEMEETINGPROVIDERTYPE OnlineMeetingProviderType = iota
    SKYPEFORBUSINESS_ONLINEMEETINGPROVIDERTYPE
    SKYPEFORCONSUMER_ONLINEMEETINGPROVIDERTYPE
    TEAMSFORBUSINESS_ONLINEMEETINGPROVIDERTYPE
)

func (i OnlineMeetingProviderType) String() string {
    return []string{"unknown", "skypeForBusiness", "skypeForConsumer", "teamsForBusiness"}[i]
}
func ParseOnlineMeetingProviderType(v string) (any, error) {
    result := UNKNOWN_ONLINEMEETINGPROVIDERTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_ONLINEMEETINGPROVIDERTYPE
        case "skypeForBusiness":
            result = SKYPEFORBUSINESS_ONLINEMEETINGPROVIDERTYPE
        case "skypeForConsumer":
            result = SKYPEFORCONSUMER_ONLINEMEETINGPROVIDERTYPE
        case "teamsForBusiness":
            result = TEAMSFORBUSINESS_ONLINEMEETINGPROVIDERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOnlineMeetingProviderType(values []OnlineMeetingProviderType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnlineMeetingProviderType) isMultiValue() bool {
    return false
}
