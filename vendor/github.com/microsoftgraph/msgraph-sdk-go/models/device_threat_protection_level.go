package models
// Device threat protection levels for the Device Threat Protection API.
type DeviceThreatProtectionLevel int

const (
    // Default Value. Do not use.
    UNAVAILABLE_DEVICETHREATPROTECTIONLEVEL DeviceThreatProtectionLevel = iota
    // Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.
    SECURED_DEVICETHREATPROTECTIONLEVEL
    // Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.
    LOW_DEVICETHREATPROTECTIONLEVEL
    // Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.
    MEDIUM_DEVICETHREATPROTECTIONLEVEL
    // Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.
    HIGH_DEVICETHREATPROTECTIONLEVEL
    // Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.
    NOTSET_DEVICETHREATPROTECTIONLEVEL
)

func (i DeviceThreatProtectionLevel) String() string {
    return []string{"unavailable", "secured", "low", "medium", "high", "notSet"}[i]
}
func ParseDeviceThreatProtectionLevel(v string) (any, error) {
    result := UNAVAILABLE_DEVICETHREATPROTECTIONLEVEL
    switch v {
        case "unavailable":
            result = UNAVAILABLE_DEVICETHREATPROTECTIONLEVEL
        case "secured":
            result = SECURED_DEVICETHREATPROTECTIONLEVEL
        case "low":
            result = LOW_DEVICETHREATPROTECTIONLEVEL
        case "medium":
            result = MEDIUM_DEVICETHREATPROTECTIONLEVEL
        case "high":
            result = HIGH_DEVICETHREATPROTECTIONLEVEL
        case "notSet":
            result = NOTSET_DEVICETHREATPROTECTIONLEVEL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceThreatProtectionLevel(values []DeviceThreatProtectionLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceThreatProtectionLevel) isMultiValue() bool {
    return false
}
