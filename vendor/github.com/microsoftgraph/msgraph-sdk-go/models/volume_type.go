package models
type VolumeType int

const (
    OPERATINGSYSTEMVOLUME_VOLUMETYPE VolumeType = iota
    FIXEDDATAVOLUME_VOLUMETYPE
    REMOVABLEDATAVOLUME_VOLUMETYPE
    UNKNOWNFUTUREVALUE_VOLUMETYPE
)

func (i VolumeType) String() string {
    return []string{"operatingSystemVolume", "fixedDataVolume", "removableDataVolume", "unknownFutureValue"}[i]
}
func ParseVolumeType(v string) (any, error) {
    result := OPERATINGSYSTEMVOLUME_VOLUMETYPE
    switch v {
        case "operatingSystemVolume":
            result = OPERATINGSYSTEMVOLUME_VOLUMETYPE
        case "fixedDataVolume":
            result = FIXEDDATAVOLUME_VOLUMETYPE
        case "removableDataVolume":
            result = REMOVABLEDATAVOLUME_VOLUMETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VOLUMETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVolumeType(values []VolumeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VolumeType) isMultiValue() bool {
    return false
}
