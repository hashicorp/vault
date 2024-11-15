package models
import (
    "math"
    "strings"
)
type WorkforceIntegrationSupportedEntities int

const (
    NONE_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 1
    SHIFT_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 2
    SWAPREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 4
    USERSHIFTPREFERENCES_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 8
    OPENSHIFT_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 16
    OPENSHIFTREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 32
    OFFERSHIFTREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 64
    UNKNOWNFUTUREVALUE_WORKFORCEINTEGRATIONSUPPORTEDENTITIES = 128
)

func (i WorkforceIntegrationSupportedEntities) String() string {
    var values []string
    options := []string{"none", "shift", "swapRequest", "userShiftPreferences", "openShift", "openShiftRequest", "offerShiftRequest", "unknownFutureValue"}
    for p := 0; p < 8; p++ {
        mantis := WorkforceIntegrationSupportedEntities(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWorkforceIntegrationSupportedEntities(v string) (any, error) {
    var result WorkforceIntegrationSupportedEntities
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "shift":
                result |= SHIFT_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "swapRequest":
                result |= SWAPREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "userShiftPreferences":
                result |= USERSHIFTPREFERENCES_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "openShift":
                result |= OPENSHIFT_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "openShiftRequest":
                result |= OPENSHIFTREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "offerShiftRequest":
                result |= OFFERSHIFTREQUEST_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_WORKFORCEINTEGRATIONSUPPORTEDENTITIES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWorkforceIntegrationSupportedEntities(values []WorkforceIntegrationSupportedEntities) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WorkforceIntegrationSupportedEntities) isMultiValue() bool {
    return true
}
