package models
type CloudPcRestorePointFrequencyType int

const (
    DEFAULT_CLOUDPCRESTOREPOINTFREQUENCYTYPE CloudPcRestorePointFrequencyType = iota
    FOURHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    SIXHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    TWELVEHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    SIXTEENHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    TWENTYFOURHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    UNKNOWNFUTUREVALUE_CLOUDPCRESTOREPOINTFREQUENCYTYPE
)

func (i CloudPcRestorePointFrequencyType) String() string {
    return []string{"default", "fourHours", "sixHours", "twelveHours", "sixteenHours", "twentyFourHours", "unknownFutureValue"}[i]
}
func ParseCloudPcRestorePointFrequencyType(v string) (any, error) {
    result := DEFAULT_CLOUDPCRESTOREPOINTFREQUENCYTYPE
    switch v {
        case "default":
            result = DEFAULT_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "fourHours":
            result = FOURHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "sixHours":
            result = SIXHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "twelveHours":
            result = TWELVEHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "sixteenHours":
            result = SIXTEENHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "twentyFourHours":
            result = TWENTYFOURHOURS_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCRESTOREPOINTFREQUENCYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcRestorePointFrequencyType(values []CloudPcRestorePointFrequencyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcRestorePointFrequencyType) isMultiValue() bool {
    return false
}
