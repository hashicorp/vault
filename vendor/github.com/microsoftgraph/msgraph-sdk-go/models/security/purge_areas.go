package security
import (
    "math"
    "strings"
)
type PurgeAreas int

const (
    MAILBOXES_PURGEAREAS = 1
    TEAMSMESSAGES_PURGEAREAS = 2
    UNKNOWNFUTUREVALUE_PURGEAREAS = 4
)

func (i PurgeAreas) String() string {
    var values []string
    options := []string{"mailboxes", "teamsMessages", "unknownFutureValue"}
    for p := 0; p < 3; p++ {
        mantis := PurgeAreas(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParsePurgeAreas(v string) (any, error) {
    var result PurgeAreas
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "mailboxes":
                result |= MAILBOXES_PURGEAREAS
            case "teamsMessages":
                result |= TEAMSMESSAGES_PURGEAREAS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_PURGEAREAS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializePurgeAreas(values []PurgeAreas) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PurgeAreas) isMultiValue() bool {
    return true
}
