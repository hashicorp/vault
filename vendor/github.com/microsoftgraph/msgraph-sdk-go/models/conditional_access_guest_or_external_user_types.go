package models
import (
    "math"
    "strings"
)
type ConditionalAccessGuestOrExternalUserTypes int

const (
    NONE_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 1
    INTERNALGUEST_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 2
    B2BCOLLABORATIONGUEST_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 4
    B2BCOLLABORATIONMEMBER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 8
    B2BDIRECTCONNECTUSER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 16
    OTHEREXTERNALUSER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 32
    SERVICEPROVIDER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 64
    UNKNOWNFUTUREVALUE_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES = 128
)

func (i ConditionalAccessGuestOrExternalUserTypes) String() string {
    var values []string
    options := []string{"none", "internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider", "unknownFutureValue"}
    for p := 0; p < 8; p++ {
        mantis := ConditionalAccessGuestOrExternalUserTypes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseConditionalAccessGuestOrExternalUserTypes(v string) (any, error) {
    var result ConditionalAccessGuestOrExternalUserTypes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "internalGuest":
                result |= INTERNALGUEST_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "b2bCollaborationGuest":
                result |= B2BCOLLABORATIONGUEST_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "b2bCollaborationMember":
                result |= B2BCOLLABORATIONMEMBER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "b2bDirectConnectUser":
                result |= B2BDIRECTCONNECTUSER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "otherExternalUser":
                result |= OTHEREXTERNALUSER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "serviceProvider":
                result |= SERVICEPROVIDER_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_CONDITIONALACCESSGUESTOREXTERNALUSERTYPES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeConditionalAccessGuestOrExternalUserTypes(values []ConditionalAccessGuestOrExternalUserTypes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConditionalAccessGuestOrExternalUserTypes) isMultiValue() bool {
    return true
}
