package models
import (
    "math"
    "strings"
)
// Type of start menu app list visibility.
type WindowsStartMenuAppListVisibilityType int

const (
    // User defined. Default value.
    USERDEFINED_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE = 1
    // Collapse the app list on the start menu.
    COLLAPSE_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE = 2
    // Removes the app list entirely from the start menu.
    REMOVE_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE = 4
    // Disables the corresponding toggle (Collapse or Remove) in the Settings app.
    DISABLESETTINGSAPP_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE = 8
)

func (i WindowsStartMenuAppListVisibilityType) String() string {
    var values []string
    options := []string{"userDefined", "collapse", "remove", "disableSettingsApp"}
    for p := 0; p < 4; p++ {
        mantis := WindowsStartMenuAppListVisibilityType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsStartMenuAppListVisibilityType(v string) (any, error) {
    var result WindowsStartMenuAppListVisibilityType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "userDefined":
                result |= USERDEFINED_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE
            case "collapse":
                result |= COLLAPSE_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE
            case "remove":
                result |= REMOVE_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE
            case "disableSettingsApp":
                result |= DISABLESETTINGSAPP_WINDOWSSTARTMENUAPPLISTVISIBILITYTYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsStartMenuAppListVisibilityType(values []WindowsStartMenuAppListVisibilityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsStartMenuAppListVisibilityType) isMultiValue() bool {
    return true
}
