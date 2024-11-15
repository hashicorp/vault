package models
// Available health states for the Device Health API
type ManagedDevicePartnerReportedHealthState int

const (
    // Device health state is not yet reported
    UNKNOWN_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE ManagedDevicePartnerReportedHealthState = iota
    // Device has been activated by a mobile threat defense partner, but has not yet reported health.
    ACTIVATED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device has been deactivated by a mobile threat defense partner. The device health is not known.
    DEACTIVATED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered secured by the mobile threat defense partner.
    SECURED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered low threat by the mobile threat defense partner.
    LOWSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered medium threat by the mobile threat defense partner.
    MEDIUMSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered high threat by the mobile threat defense partner.
    HIGHSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered unresponsive by the mobile threat defense partner. The device health is not known.
    UNRESPONSIVE_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered compromised by the Threat Defense partner. This means the device has an active Threat or Risk which cannot be easily remediated by the end user and the user should contact their IT Admin.
    COMPROMISED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    // Device is considered misconfigured with the Threat Defense partner. This means the device is missing a required profile or configuration for the Threat Defense Partner to function properly and is thus threat or risk analysis is not able to complete.
    MISCONFIGURED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
)

func (i ManagedDevicePartnerReportedHealthState) String() string {
    return []string{"unknown", "activated", "deactivated", "secured", "lowSeverity", "mediumSeverity", "highSeverity", "unresponsive", "compromised", "misconfigured"}[i]
}
func ParseManagedDevicePartnerReportedHealthState(v string) (any, error) {
    result := UNKNOWN_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "activated":
            result = ACTIVATED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "deactivated":
            result = DEACTIVATED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "secured":
            result = SECURED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "lowSeverity":
            result = LOWSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "mediumSeverity":
            result = MEDIUMSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "highSeverity":
            result = HIGHSEVERITY_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "unresponsive":
            result = UNRESPONSIVE_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "compromised":
            result = COMPROMISED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        case "misconfigured":
            result = MISCONFIGURED_MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedDevicePartnerReportedHealthState(values []ManagedDevicePartnerReportedHealthState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedDevicePartnerReportedHealthState) isMultiValue() bool {
    return false
}
