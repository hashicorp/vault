package models
// Type of display modes for the start menu.
type WindowsStartMenuModeType int

const (
    // User defined. Default value.
    USERDEFINED_WINDOWSSTARTMENUMODETYPE WindowsStartMenuModeType = iota
    // Full screen.
    FULLSCREEN_WINDOWSSTARTMENUMODETYPE
    // Non-full screen.
    NONFULLSCREEN_WINDOWSSTARTMENUMODETYPE
)

func (i WindowsStartMenuModeType) String() string {
    return []string{"userDefined", "fullScreen", "nonFullScreen"}[i]
}
func ParseWindowsStartMenuModeType(v string) (any, error) {
    result := USERDEFINED_WINDOWSSTARTMENUMODETYPE
    switch v {
        case "userDefined":
            result = USERDEFINED_WINDOWSSTARTMENUMODETYPE
        case "fullScreen":
            result = FULLSCREEN_WINDOWSSTARTMENUMODETYPE
        case "nonFullScreen":
            result = NONFULLSCREEN_WINDOWSSTARTMENUMODETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsStartMenuModeType(values []WindowsStartMenuModeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsStartMenuModeType) isMultiValue() bool {
    return false
}
