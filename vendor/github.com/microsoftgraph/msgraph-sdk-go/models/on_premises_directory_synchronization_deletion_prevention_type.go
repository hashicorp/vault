package models
type OnPremisesDirectorySynchronizationDeletionPreventionType int

const (
    DISABLED_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE OnPremisesDirectorySynchronizationDeletionPreventionType = iota
    ENABLEDFORCOUNT_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
    ENABLEDFORPERCENTAGE_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
    UNKNOWNFUTUREVALUE_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
)

func (i OnPremisesDirectorySynchronizationDeletionPreventionType) String() string {
    return []string{"disabled", "enabledForCount", "enabledForPercentage", "unknownFutureValue"}[i]
}
func ParseOnPremisesDirectorySynchronizationDeletionPreventionType(v string) (any, error) {
    result := DISABLED_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
    switch v {
        case "disabled":
            result = DISABLED_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
        case "enabledForCount":
            result = ENABLEDFORCOUNT_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
        case "enabledForPercentage":
            result = ENABLEDFORPERCENTAGE_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ONPREMISESDIRECTORYSYNCHRONIZATIONDELETIONPREVENTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOnPremisesDirectorySynchronizationDeletionPreventionType(values []OnPremisesDirectorySynchronizationDeletionPreventionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OnPremisesDirectorySynchronizationDeletionPreventionType) isMultiValue() bool {
    return false
}
