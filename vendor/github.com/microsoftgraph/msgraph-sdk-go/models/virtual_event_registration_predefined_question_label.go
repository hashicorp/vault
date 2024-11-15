package models
type VirtualEventRegistrationPredefinedQuestionLabel int

const (
    STREET_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL VirtualEventRegistrationPredefinedQuestionLabel = iota
    CITY_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    STATE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    POSTALCODE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    COUNTRYORREGION_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    INDUSTRY_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    JOBTITLE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    ORGANIZATION_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    UNKNOWNFUTUREVALUE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
)

func (i VirtualEventRegistrationPredefinedQuestionLabel) String() string {
    return []string{"street", "city", "state", "postalCode", "countryOrRegion", "industry", "jobTitle", "organization", "unknownFutureValue"}[i]
}
func ParseVirtualEventRegistrationPredefinedQuestionLabel(v string) (any, error) {
    result := STREET_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
    switch v {
        case "street":
            result = STREET_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "city":
            result = CITY_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "state":
            result = STATE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "postalCode":
            result = POSTALCODE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "countryOrRegion":
            result = COUNTRYORREGION_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "industry":
            result = INDUSTRY_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "jobTitle":
            result = JOBTITLE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "organization":
            result = ORGANIZATION_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIRTUALEVENTREGISTRATIONPREDEFINEDQUESTIONLABEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVirtualEventRegistrationPredefinedQuestionLabel(values []VirtualEventRegistrationPredefinedQuestionLabel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VirtualEventRegistrationPredefinedQuestionLabel) isMultiValue() bool {
    return false
}
