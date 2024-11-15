package models
type PhoneType int

const (
    HOME_PHONETYPE PhoneType = iota
    BUSINESS_PHONETYPE
    MOBILE_PHONETYPE
    OTHER_PHONETYPE
    ASSISTANT_PHONETYPE
    HOMEFAX_PHONETYPE
    BUSINESSFAX_PHONETYPE
    OTHERFAX_PHONETYPE
    PAGER_PHONETYPE
    RADIO_PHONETYPE
)

func (i PhoneType) String() string {
    return []string{"home", "business", "mobile", "other", "assistant", "homeFax", "businessFax", "otherFax", "pager", "radio"}[i]
}
func ParsePhoneType(v string) (any, error) {
    result := HOME_PHONETYPE
    switch v {
        case "home":
            result = HOME_PHONETYPE
        case "business":
            result = BUSINESS_PHONETYPE
        case "mobile":
            result = MOBILE_PHONETYPE
        case "other":
            result = OTHER_PHONETYPE
        case "assistant":
            result = ASSISTANT_PHONETYPE
        case "homeFax":
            result = HOMEFAX_PHONETYPE
        case "businessFax":
            result = BUSINESSFAX_PHONETYPE
        case "otherFax":
            result = OTHERFAX_PHONETYPE
        case "pager":
            result = PAGER_PHONETYPE
        case "radio":
            result = RADIO_PHONETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePhoneType(values []PhoneType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PhoneType) isMultiValue() bool {
    return false
}
