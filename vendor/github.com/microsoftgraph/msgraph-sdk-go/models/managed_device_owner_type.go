package models
// Owner type of device.
type ManagedDeviceOwnerType int

const (
    // Unknown.
    UNKNOWN_MANAGEDDEVICEOWNERTYPE ManagedDeviceOwnerType = iota
    // Owned by company.
    COMPANY_MANAGEDDEVICEOWNERTYPE
    // Owned by person.
    PERSONAL_MANAGEDDEVICEOWNERTYPE
)

func (i ManagedDeviceOwnerType) String() string {
    return []string{"unknown", "company", "personal"}[i]
}
func ParseManagedDeviceOwnerType(v string) (any, error) {
    result := UNKNOWN_MANAGEDDEVICEOWNERTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_MANAGEDDEVICEOWNERTYPE
        case "company":
            result = COMPANY_MANAGEDDEVICEOWNERTYPE
        case "personal":
            result = PERSONAL_MANAGEDDEVICEOWNERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedDeviceOwnerType(values []ManagedDeviceOwnerType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedDeviceOwnerType) isMultiValue() bool {
    return false
}
