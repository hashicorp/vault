package models
type DiskType int

const (
    // Enum member for unknown or default diskType.
    UNKNOWN_DISKTYPE DiskType = iota
    // Enum member for HDD devices.
    HDD_DISKTYPE
    // Enum member for SSD devices.
    SSD_DISKTYPE
    // Evolvable enumeration sentinel value.Do not use.
    UNKNOWNFUTUREVALUE_DISKTYPE
)

func (i DiskType) String() string {
    return []string{"unknown", "hdd", "ssd", "unknownFutureValue"}[i]
}
func ParseDiskType(v string) (any, error) {
    result := UNKNOWN_DISKTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_DISKTYPE
        case "hdd":
            result = HDD_DISKTYPE
        case "ssd":
            result = SSD_DISKTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DISKTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDiskType(values []DiskType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DiskType) isMultiValue() bool {
    return false
}
