package models
import (
    "math"
    "strings"
)
type NativeAuthenticationApisEnabled int

const (
    NONE_NATIVEAUTHENTICATIONAPISENABLED = 1
    ALL_NATIVEAUTHENTICATIONAPISENABLED = 2
    UNKNOWNFUTUREVALUE_NATIVEAUTHENTICATIONAPISENABLED = 4
)

func (i NativeAuthenticationApisEnabled) String() string {
    var values []string
    options := []string{"none", "all", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := NativeAuthenticationApisEnabled(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseNativeAuthenticationApisEnabled(v string) (any, error) {
    var result NativeAuthenticationApisEnabled
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_NATIVEAUTHENTICATIONAPISENABLED
            case "all":
                result |= ALL_NATIVEAUTHENTICATIONAPISENABLED
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_NATIVEAUTHENTICATIONAPISENABLED
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeNativeAuthenticationApisEnabled(values []NativeAuthenticationApisEnabled) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i NativeAuthenticationApisEnabled) isMultiValue() bool {
    return true
}
