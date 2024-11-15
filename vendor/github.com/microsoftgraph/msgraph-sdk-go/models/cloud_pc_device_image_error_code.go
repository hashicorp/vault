package models
type CloudPcDeviceImageErrorCode int

const (
    INTERNALSERVERERROR_CLOUDPCDEVICEIMAGEERRORCODE CloudPcDeviceImageErrorCode = iota
    SOURCEIMAGENOTFOUND_CLOUDPCDEVICEIMAGEERRORCODE
    OSVERSIONNOTSUPPORTED_CLOUDPCDEVICEIMAGEERRORCODE
    SOURCEIMAGEINVALID_CLOUDPCDEVICEIMAGEERRORCODE
    SOURCEIMAGENOTGENERALIZED_CLOUDPCDEVICEIMAGEERRORCODE
    UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGEERRORCODE
    VMALREADYAZUREADJOINED_CLOUDPCDEVICEIMAGEERRORCODE
    PAIDSOURCEIMAGENOTSUPPORT_CLOUDPCDEVICEIMAGEERRORCODE
    SOURCEIMAGENOTSUPPORTCUSTOMIZEVMNAME_CLOUDPCDEVICEIMAGEERRORCODE
    SOURCEIMAGESIZEEXCEEDSLIMITATION_CLOUDPCDEVICEIMAGEERRORCODE
)

func (i CloudPcDeviceImageErrorCode) String() string {
    return []string{"internalServerError", "sourceImageNotFound", "osVersionNotSupported", "sourceImageInvalid", "sourceImageNotGeneralized", "unknownFutureValue", "vmAlreadyAzureAdjoined", "paidSourceImageNotSupport", "sourceImageNotSupportCustomizeVMName", "sourceImageSizeExceedsLimitation"}[i]
}
func ParseCloudPcDeviceImageErrorCode(v string) (any, error) {
    result := INTERNALSERVERERROR_CLOUDPCDEVICEIMAGEERRORCODE
    switch v {
        case "internalServerError":
            result = INTERNALSERVERERROR_CLOUDPCDEVICEIMAGEERRORCODE
        case "sourceImageNotFound":
            result = SOURCEIMAGENOTFOUND_CLOUDPCDEVICEIMAGEERRORCODE
        case "osVersionNotSupported":
            result = OSVERSIONNOTSUPPORTED_CLOUDPCDEVICEIMAGEERRORCODE
        case "sourceImageInvalid":
            result = SOURCEIMAGEINVALID_CLOUDPCDEVICEIMAGEERRORCODE
        case "sourceImageNotGeneralized":
            result = SOURCEIMAGENOTGENERALIZED_CLOUDPCDEVICEIMAGEERRORCODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGEERRORCODE
        case "vmAlreadyAzureAdjoined":
            result = VMALREADYAZUREADJOINED_CLOUDPCDEVICEIMAGEERRORCODE
        case "paidSourceImageNotSupport":
            result = PAIDSOURCEIMAGENOTSUPPORT_CLOUDPCDEVICEIMAGEERRORCODE
        case "sourceImageNotSupportCustomizeVMName":
            result = SOURCEIMAGENOTSUPPORTCUSTOMIZEVMNAME_CLOUDPCDEVICEIMAGEERRORCODE
        case "sourceImageSizeExceedsLimitation":
            result = SOURCEIMAGESIZEEXCEEDSLIMITATION_CLOUDPCDEVICEIMAGEERRORCODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcDeviceImageErrorCode(values []CloudPcDeviceImageErrorCode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcDeviceImageErrorCode) isMultiValue() bool {
    return false
}
