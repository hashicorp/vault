package models
import (
    "math"
    "strings"
)
// Scheduled the update installation on the weeks of the month
type WindowsUpdateForBusinessUpdateWeeks int

const (
    // Allow the user to set.
    USERDEFINED_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 1
    // Scheduled the update installation on the first week of the month
    FIRSTWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 2
    // Scheduled the update installation on the second week of the month
    SECONDWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 4
    // Scheduled the update installation on the third week of the month
    THIRDWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 8
    // Scheduled the update installation on the fourth week of the month
    FOURTHWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 16
    // Scheduled the update installation on every week of the month
    EVERYWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 32
    // Evolvable enum member
    UNKNOWNFUTUREVALUE_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS = 64
)

func (i WindowsUpdateForBusinessUpdateWeeks) String() string {
    var values []string
    options := []string{"userDefined", "firstWeek", "secondWeek", "thirdWeek", "fourthWeek", "everyWeek", "unknownFutureValue"}
    for p := 0; p < 7; p++ {
        mantis := WindowsUpdateForBusinessUpdateWeeks(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsUpdateForBusinessUpdateWeeks(v string) (any, error) {
    var result WindowsUpdateForBusinessUpdateWeeks
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "userDefined":
                result |= USERDEFINED_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "firstWeek":
                result |= FIRSTWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "secondWeek":
                result |= SECONDWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "thirdWeek":
                result |= THIRDWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "fourthWeek":
                result |= FOURTHWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "everyWeek":
                result |= EVERYWEEK_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_WINDOWSUPDATEFORBUSINESSUPDATEWEEKS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsUpdateForBusinessUpdateWeeks(values []WindowsUpdateForBusinessUpdateWeeks) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsUpdateForBusinessUpdateWeeks) isMultiValue() bool {
    return true
}
