package models
type MicrosoftStoreForBusinessLicenseType int

const (
    OFFLINE_MICROSOFTSTOREFORBUSINESSLICENSETYPE MicrosoftStoreForBusinessLicenseType = iota
    ONLINE_MICROSOFTSTOREFORBUSINESSLICENSETYPE
)

func (i MicrosoftStoreForBusinessLicenseType) String() string {
    return []string{"offline", "online"}[i]
}
func ParseMicrosoftStoreForBusinessLicenseType(v string) (any, error) {
    result := OFFLINE_MICROSOFTSTOREFORBUSINESSLICENSETYPE
    switch v {
        case "offline":
            result = OFFLINE_MICROSOFTSTOREFORBUSINESSLICENSETYPE
        case "online":
            result = ONLINE_MICROSOFTSTOREFORBUSINESSLICENSETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMicrosoftStoreForBusinessLicenseType(values []MicrosoftStoreForBusinessLicenseType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MicrosoftStoreForBusinessLicenseType) isMultiValue() bool {
    return false
}
