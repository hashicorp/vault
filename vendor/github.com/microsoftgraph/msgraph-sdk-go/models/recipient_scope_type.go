package models
import (
    "math"
    "strings"
)
type RecipientScopeType int

const (
    NONE_RECIPIENTSCOPETYPE = 1
    INTERNAL_RECIPIENTSCOPETYPE = 2
    EXTERNAL_RECIPIENTSCOPETYPE = 4
    EXTERNALPARTNER_RECIPIENTSCOPETYPE = 8
    EXTERNALNONPARTNER_RECIPIENTSCOPETYPE = 16
)

func (i RecipientScopeType) String() string {
    var values []string
    options := []string{"none", "internal", "external", "externalPartner", "externalNonPartner"}
    for p := 0; p < 5; p++ {
        mantis := RecipientScopeType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseRecipientScopeType(v string) (any, error) {
    var result RecipientScopeType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_RECIPIENTSCOPETYPE
            case "internal":
                result |= INTERNAL_RECIPIENTSCOPETYPE
            case "external":
                result |= EXTERNAL_RECIPIENTSCOPETYPE
            case "externalPartner":
                result |= EXTERNALPARTNER_RECIPIENTSCOPETYPE
            case "externalNonPartner":
                result |= EXTERNALNONPARTNER_RECIPIENTSCOPETYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeRecipientScopeType(values []RecipientScopeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RecipientScopeType) isMultiValue() bool {
    return true
}
