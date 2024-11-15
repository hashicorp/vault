package models
// Work From Anywhere windows device upgrade eligibility status.
type OperatingSystemUpgradeEligibility int

const (
    // The device is upgraded to latest version of windows.
    UPGRADED_OPERATINGSYSTEMUPGRADEELIGIBILITY OperatingSystemUpgradeEligibility = iota
    // Not enough data available to compute the eligibility of device for windows upgrade.
    UNKNOWN_OPERATINGSYSTEMUPGRADEELIGIBILITY
    // The device is not capable for windows upgrade.
    NOTCAPABLE_OPERATINGSYSTEMUPGRADEELIGIBILITY
    // The device is capable for windows upgrade.
    CAPABLE_OPERATINGSYSTEMUPGRADEELIGIBILITY
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_OPERATINGSYSTEMUPGRADEELIGIBILITY
)

func (i OperatingSystemUpgradeEligibility) String() string {
    return []string{"upgraded", "unknown", "notCapable", "capable", "unknownFutureValue"}[i]
}
func ParseOperatingSystemUpgradeEligibility(v string) (any, error) {
    result := UPGRADED_OPERATINGSYSTEMUPGRADEELIGIBILITY
    switch v {
        case "upgraded":
            result = UPGRADED_OPERATINGSYSTEMUPGRADEELIGIBILITY
        case "unknown":
            result = UNKNOWN_OPERATINGSYSTEMUPGRADEELIGIBILITY
        case "notCapable":
            result = NOTCAPABLE_OPERATINGSYSTEMUPGRADEELIGIBILITY
        case "capable":
            result = CAPABLE_OPERATINGSYSTEMUPGRADEELIGIBILITY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_OPERATINGSYSTEMUPGRADEELIGIBILITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOperatingSystemUpgradeEligibility(values []OperatingSystemUpgradeEligibility) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OperatingSystemUpgradeEligibility) isMultiValue() bool {
    return false
}
