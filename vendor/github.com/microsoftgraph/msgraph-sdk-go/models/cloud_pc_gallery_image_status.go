package models
type CloudPcGalleryImageStatus int

const (
    SUPPORTED_CLOUDPCGALLERYIMAGESTATUS CloudPcGalleryImageStatus = iota
    SUPPORTEDWITHWARNING_CLOUDPCGALLERYIMAGESTATUS
    NOTSUPPORTED_CLOUDPCGALLERYIMAGESTATUS
    UNKNOWNFUTUREVALUE_CLOUDPCGALLERYIMAGESTATUS
)

func (i CloudPcGalleryImageStatus) String() string {
    return []string{"supported", "supportedWithWarning", "notSupported", "unknownFutureValue"}[i]
}
func ParseCloudPcGalleryImageStatus(v string) (any, error) {
    result := SUPPORTED_CLOUDPCGALLERYIMAGESTATUS
    switch v {
        case "supported":
            result = SUPPORTED_CLOUDPCGALLERYIMAGESTATUS
        case "supportedWithWarning":
            result = SUPPORTEDWITHWARNING_CLOUDPCGALLERYIMAGESTATUS
        case "notSupported":
            result = NOTSUPPORTED_CLOUDPCGALLERYIMAGESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCGALLERYIMAGESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcGalleryImageStatus(values []CloudPcGalleryImageStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcGalleryImageStatus) isMultiValue() bool {
    return false
}
