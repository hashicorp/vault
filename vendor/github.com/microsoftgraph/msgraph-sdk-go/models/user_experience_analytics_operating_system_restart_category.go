package models
// Operating System restart category.
type UserExperienceAnalyticsOperatingSystemRestartCategory int

const (
    // Default. Set to unknown if device operating system restart category has not yet been calculated.
    UNKNOWN_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY UserExperienceAnalyticsOperatingSystemRestartCategory = iota
    // Indicates that the device operating system restart is along with an update.
    RESTARTWITHUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is without update.
    RESTARTWITHOUTUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is due to a specific stop error.
    BLUESCREEN_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is due to shutdown with update.
    SHUTDOWNWITHUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is due to shutdown without update.
    SHUTDOWNWITHOUTUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is due to update long power-button press.
    LONGPOWERBUTTONPRESS_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restart is due to boot error.
    BOOTERROR_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Indicates that the device operating system restarted after an update.
    UPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
)

func (i UserExperienceAnalyticsOperatingSystemRestartCategory) String() string {
    return []string{"unknown", "restartWithUpdate", "restartWithoutUpdate", "blueScreen", "shutdownWithUpdate", "shutdownWithoutUpdate", "longPowerButtonPress", "bootError", "update", "unknownFutureValue"}[i]
}
func ParseUserExperienceAnalyticsOperatingSystemRestartCategory(v string) (any, error) {
    result := UNKNOWN_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
    switch v {
        case "unknown":
            result = UNKNOWN_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "restartWithUpdate":
            result = RESTARTWITHUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "restartWithoutUpdate":
            result = RESTARTWITHOUTUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "blueScreen":
            result = BLUESCREEN_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "shutdownWithUpdate":
            result = SHUTDOWNWITHUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "shutdownWithoutUpdate":
            result = SHUTDOWNWITHOUTUPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "longPowerButtonPress":
            result = LONGPOWERBUTTONPRESS_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "bootError":
            result = BOOTERROR_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "update":
            result = UPDATE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USEREXPERIENCEANALYTICSOPERATINGSYSTEMRESTARTCATEGORY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserExperienceAnalyticsOperatingSystemRestartCategory(values []UserExperienceAnalyticsOperatingSystemRestartCategory) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserExperienceAnalyticsOperatingSystemRestartCategory) isMultiValue() bool {
    return false
}
