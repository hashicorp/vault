package models
type LocationType int

const (
    DEFAULT_LOCATIONTYPE LocationType = iota
    CONFERENCEROOM_LOCATIONTYPE
    HOMEADDRESS_LOCATIONTYPE
    BUSINESSADDRESS_LOCATIONTYPE
    GEOCOORDINATES_LOCATIONTYPE
    STREETADDRESS_LOCATIONTYPE
    HOTEL_LOCATIONTYPE
    RESTAURANT_LOCATIONTYPE
    LOCALBUSINESS_LOCATIONTYPE
    POSTALADDRESS_LOCATIONTYPE
)

func (i LocationType) String() string {
    return []string{"default", "conferenceRoom", "homeAddress", "businessAddress", "geoCoordinates", "streetAddress", "hotel", "restaurant", "localBusiness", "postalAddress"}[i]
}
func ParseLocationType(v string) (any, error) {
    result := DEFAULT_LOCATIONTYPE
    switch v {
        case "default":
            result = DEFAULT_LOCATIONTYPE
        case "conferenceRoom":
            result = CONFERENCEROOM_LOCATIONTYPE
        case "homeAddress":
            result = HOMEADDRESS_LOCATIONTYPE
        case "businessAddress":
            result = BUSINESSADDRESS_LOCATIONTYPE
        case "geoCoordinates":
            result = GEOCOORDINATES_LOCATIONTYPE
        case "streetAddress":
            result = STREETADDRESS_LOCATIONTYPE
        case "hotel":
            result = HOTEL_LOCATIONTYPE
        case "restaurant":
            result = RESTAURANT_LOCATIONTYPE
        case "localBusiness":
            result = LOCALBUSINESS_LOCATIONTYPE
        case "postalAddress":
            result = POSTALADDRESS_LOCATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeLocationType(values []LocationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i LocationType) isMultiValue() bool {
    return false
}
