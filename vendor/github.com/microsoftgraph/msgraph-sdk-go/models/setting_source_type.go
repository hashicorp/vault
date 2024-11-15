package models
type SettingSourceType int

const (
    DEVICECONFIGURATION_SETTINGSOURCETYPE SettingSourceType = iota
    DEVICEINTENT_SETTINGSOURCETYPE
)

func (i SettingSourceType) String() string {
    return []string{"deviceConfiguration", "deviceIntent"}[i]
}
func ParseSettingSourceType(v string) (any, error) {
    result := DEVICECONFIGURATION_SETTINGSOURCETYPE
    switch v {
        case "deviceConfiguration":
            result = DEVICECONFIGURATION_SETTINGSOURCETYPE
        case "deviceIntent":
            result = DEVICEINTENT_SETTINGSOURCETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSettingSourceType(values []SettingSourceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SettingSourceType) isMultiValue() bool {
    return false
}
