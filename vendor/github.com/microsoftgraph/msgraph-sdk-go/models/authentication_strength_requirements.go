package models
import (
    "math"
    "strings"
)
type AuthenticationStrengthRequirements int

const (
    NONE_AUTHENTICATIONSTRENGTHREQUIREMENTS = 1
    MFA_AUTHENTICATIONSTRENGTHREQUIREMENTS = 2
    UNKNOWNFUTUREVALUE_AUTHENTICATIONSTRENGTHREQUIREMENTS = 4
)

func (i AuthenticationStrengthRequirements) String() string {
    var values []string
    options := []string{"none", "mfa", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := AuthenticationStrengthRequirements(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseAuthenticationStrengthRequirements(v string) (any, error) {
    var result AuthenticationStrengthRequirements
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_AUTHENTICATIONSTRENGTHREQUIREMENTS
            case "mfa":
                result |= MFA_AUTHENTICATIONSTRENGTHREQUIREMENTS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_AUTHENTICATIONSTRENGTHREQUIREMENTS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeAuthenticationStrengthRequirements(values []AuthenticationStrengthRequirements) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationStrengthRequirements) isMultiValue() bool {
    return true
}
