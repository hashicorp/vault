package models
type PayloadBrand int

const (
    UNKNOWN_PAYLOADBRAND PayloadBrand = iota
    OTHER_PAYLOADBRAND
    AMERICANEXPRESS_PAYLOADBRAND
    CAPITALONE_PAYLOADBRAND
    DHL_PAYLOADBRAND
    DOCUSIGN_PAYLOADBRAND
    DROPBOX_PAYLOADBRAND
    FACEBOOK_PAYLOADBRAND
    FIRSTAMERICAN_PAYLOADBRAND
    MICROSOFT_PAYLOADBRAND
    NETFLIX_PAYLOADBRAND
    SCOTIABANK_PAYLOADBRAND
    SENDGRID_PAYLOADBRAND
    STEWARTTITLE_PAYLOADBRAND
    TESCO_PAYLOADBRAND
    WELLSFARGO_PAYLOADBRAND
    SYRINXCLOUD_PAYLOADBRAND
    ADOBE_PAYLOADBRAND
    TEAMS_PAYLOADBRAND
    ZOOM_PAYLOADBRAND
    UNKNOWNFUTUREVALUE_PAYLOADBRAND
)

func (i PayloadBrand) String() string {
    return []string{"unknown", "other", "americanExpress", "capitalOne", "dhl", "docuSign", "dropbox", "facebook", "firstAmerican", "microsoft", "netflix", "scotiabank", "sendGrid", "stewartTitle", "tesco", "wellsFargo", "syrinxCloud", "adobe", "teams", "zoom", "unknownFutureValue"}[i]
}
func ParsePayloadBrand(v string) (any, error) {
    result := UNKNOWN_PAYLOADBRAND
    switch v {
        case "unknown":
            result = UNKNOWN_PAYLOADBRAND
        case "other":
            result = OTHER_PAYLOADBRAND
        case "americanExpress":
            result = AMERICANEXPRESS_PAYLOADBRAND
        case "capitalOne":
            result = CAPITALONE_PAYLOADBRAND
        case "dhl":
            result = DHL_PAYLOADBRAND
        case "docuSign":
            result = DOCUSIGN_PAYLOADBRAND
        case "dropbox":
            result = DROPBOX_PAYLOADBRAND
        case "facebook":
            result = FACEBOOK_PAYLOADBRAND
        case "firstAmerican":
            result = FIRSTAMERICAN_PAYLOADBRAND
        case "microsoft":
            result = MICROSOFT_PAYLOADBRAND
        case "netflix":
            result = NETFLIX_PAYLOADBRAND
        case "scotiabank":
            result = SCOTIABANK_PAYLOADBRAND
        case "sendGrid":
            result = SENDGRID_PAYLOADBRAND
        case "stewartTitle":
            result = STEWARTTITLE_PAYLOADBRAND
        case "tesco":
            result = TESCO_PAYLOADBRAND
        case "wellsFargo":
            result = WELLSFARGO_PAYLOADBRAND
        case "syrinxCloud":
            result = SYRINXCLOUD_PAYLOADBRAND
        case "adobe":
            result = ADOBE_PAYLOADBRAND
        case "teams":
            result = TEAMS_PAYLOADBRAND
        case "zoom":
            result = ZOOM_PAYLOADBRAND
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PAYLOADBRAND
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePayloadBrand(values []PayloadBrand) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PayloadBrand) isMultiValue() bool {
    return false
}
