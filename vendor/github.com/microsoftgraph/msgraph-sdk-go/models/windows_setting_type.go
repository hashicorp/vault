package models
type WindowsSettingType int

const (
    ROAMING_WINDOWSSETTINGTYPE WindowsSettingType = iota
    BACKUP_WINDOWSSETTINGTYPE
    UNKNOWNFUTUREVALUE_WINDOWSSETTINGTYPE
)

func (i WindowsSettingType) String() string {
    return []string{"roaming", "backup", "unknownFutureValue"}[i]
}
func ParseWindowsSettingType(v string) (any, error) {
    result := ROAMING_WINDOWSSETTINGTYPE
    switch v {
        case "roaming":
            result = ROAMING_WINDOWSSETTINGTYPE
        case "backup":
            result = BACKUP_WINDOWSSETTINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WINDOWSSETTINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsSettingType(values []WindowsSettingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsSettingType) isMultiValue() bool {
    return false
}
