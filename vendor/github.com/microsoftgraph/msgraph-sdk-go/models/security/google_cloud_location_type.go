package security
type GoogleCloudLocationType int

const (
    UNKNOWN_GOOGLECLOUDLOCATIONTYPE GoogleCloudLocationType = iota
    REGIONAL_GOOGLECLOUDLOCATIONTYPE
    ZONAL_GOOGLECLOUDLOCATIONTYPE
    GLOBAL_GOOGLECLOUDLOCATIONTYPE
    UNKNOWNFUTUREVALUE_GOOGLECLOUDLOCATIONTYPE
)

func (i GoogleCloudLocationType) String() string {
    return []string{"unknown", "regional", "zonal", "global", "unknownFutureValue"}[i]
}
func ParseGoogleCloudLocationType(v string) (any, error) {
    result := UNKNOWN_GOOGLECLOUDLOCATIONTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_GOOGLECLOUDLOCATIONTYPE
        case "regional":
            result = REGIONAL_GOOGLECLOUDLOCATIONTYPE
        case "zonal":
            result = ZONAL_GOOGLECLOUDLOCATIONTYPE
        case "global":
            result = GLOBAL_GOOGLECLOUDLOCATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_GOOGLECLOUDLOCATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeGoogleCloudLocationType(values []GoogleCloudLocationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i GoogleCloudLocationType) isMultiValue() bool {
    return false
}
