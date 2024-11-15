package models
import (
    "math"
    "strings"
)
type ClonableTeamParts int

const (
    APPS_CLONABLETEAMPARTS = 1
    TABS_CLONABLETEAMPARTS = 2
    SETTINGS_CLONABLETEAMPARTS = 4
    CHANNELS_CLONABLETEAMPARTS = 8
    MEMBERS_CLONABLETEAMPARTS = 16
)

func (i ClonableTeamParts) String() string {
    var values []string
    options := []string{"apps", "tabs", "settings", "channels", "members"}
    for p := 0; p < 5; p++ {
        mantis := ClonableTeamParts(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseClonableTeamParts(v string) (any, error) {
    var result ClonableTeamParts
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "apps":
                result |= APPS_CLONABLETEAMPARTS
            case "tabs":
                result |= TABS_CLONABLETEAMPARTS
            case "settings":
                result |= SETTINGS_CLONABLETEAMPARTS
            case "channels":
                result |= CHANNELS_CLONABLETEAMPARTS
            case "members":
                result |= MEMBERS_CLONABLETEAMPARTS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeClonableTeamParts(values []ClonableTeamParts) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ClonableTeamParts) isMultiValue() bool {
    return true
}
