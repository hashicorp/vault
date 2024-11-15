package models
type MicrosoftManagedDesktopType int

const (
    NOTMANAGED_MICROSOFTMANAGEDDESKTOPTYPE MicrosoftManagedDesktopType = iota
    PREMIUMMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
    STANDARDMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
    STARTERMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
    UNKNOWNFUTUREVALUE_MICROSOFTMANAGEDDESKTOPTYPE
)

func (i MicrosoftManagedDesktopType) String() string {
    return []string{"notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"}[i]
}
func ParseMicrosoftManagedDesktopType(v string) (any, error) {
    result := NOTMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
    switch v {
        case "notManaged":
            result = NOTMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
        case "premiumManaged":
            result = PREMIUMMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
        case "standardManaged":
            result = STANDARDMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
        case "starterManaged":
            result = STARTERMANAGED_MICROSOFTMANAGEDDESKTOPTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MICROSOFTMANAGEDDESKTOPTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMicrosoftManagedDesktopType(values []MicrosoftManagedDesktopType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MicrosoftManagedDesktopType) isMultiValue() bool {
    return false
}
