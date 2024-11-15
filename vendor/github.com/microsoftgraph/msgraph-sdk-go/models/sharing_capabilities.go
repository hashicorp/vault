package models
type SharingCapabilities int

const (
    DISABLED_SHARINGCAPABILITIES SharingCapabilities = iota
    EXTERNALUSERSHARINGONLY_SHARINGCAPABILITIES
    EXTERNALUSERANDGUESTSHARING_SHARINGCAPABILITIES
    EXISTINGEXTERNALUSERSHARINGONLY_SHARINGCAPABILITIES
    UNKNOWNFUTUREVALUE_SHARINGCAPABILITIES
)

func (i SharingCapabilities) String() string {
    return []string{"disabled", "externalUserSharingOnly", "externalUserAndGuestSharing", "existingExternalUserSharingOnly", "unknownFutureValue"}[i]
}
func ParseSharingCapabilities(v string) (any, error) {
    result := DISABLED_SHARINGCAPABILITIES
    switch v {
        case "disabled":
            result = DISABLED_SHARINGCAPABILITIES
        case "externalUserSharingOnly":
            result = EXTERNALUSERSHARINGONLY_SHARINGCAPABILITIES
        case "externalUserAndGuestSharing":
            result = EXTERNALUSERANDGUESTSHARING_SHARINGCAPABILITIES
        case "existingExternalUserSharingOnly":
            result = EXISTINGEXTERNALUSERSHARINGONLY_SHARINGCAPABILITIES
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SHARINGCAPABILITIES
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSharingCapabilities(values []SharingCapabilities) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SharingCapabilities) isMultiValue() bool {
    return false
}
