package models
// Device registration status.
type DeviceRegistrationState int

const (
    // The device is not registered.
    NOTREGISTERED_DEVICEREGISTRATIONSTATE DeviceRegistrationState = iota
    // The device is registered.
    REGISTERED_DEVICEREGISTRATIONSTATE
    // The device has been blocked, wiped or retired.
    REVOKED_DEVICEREGISTRATIONSTATE
    // The device has a key conflict.
    KEYCONFLICT_DEVICEREGISTRATIONSTATE
    // The device is pending approval.
    APPROVALPENDING_DEVICEREGISTRATIONSTATE
    // The device certificate has been reset.
    CERTIFICATERESET_DEVICEREGISTRATIONSTATE
    // The device is not registered and pending enrollment.
    NOTREGISTEREDPENDINGENROLLMENT_DEVICEREGISTRATIONSTATE
    // The device registration status is unknown.
    UNKNOWN_DEVICEREGISTRATIONSTATE
)

func (i DeviceRegistrationState) String() string {
    return []string{"notRegistered", "registered", "revoked", "keyConflict", "approvalPending", "certificateReset", "notRegisteredPendingEnrollment", "unknown"}[i]
}
func ParseDeviceRegistrationState(v string) (any, error) {
    result := NOTREGISTERED_DEVICEREGISTRATIONSTATE
    switch v {
        case "notRegistered":
            result = NOTREGISTERED_DEVICEREGISTRATIONSTATE
        case "registered":
            result = REGISTERED_DEVICEREGISTRATIONSTATE
        case "revoked":
            result = REVOKED_DEVICEREGISTRATIONSTATE
        case "keyConflict":
            result = KEYCONFLICT_DEVICEREGISTRATIONSTATE
        case "approvalPending":
            result = APPROVALPENDING_DEVICEREGISTRATIONSTATE
        case "certificateReset":
            result = CERTIFICATERESET_DEVICEREGISTRATIONSTATE
        case "notRegisteredPendingEnrollment":
            result = NOTREGISTEREDPENDINGENROLLMENT_DEVICEREGISTRATIONSTATE
        case "unknown":
            result = UNKNOWN_DEVICEREGISTRATIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceRegistrationState(values []DeviceRegistrationState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceRegistrationState) isMultiValue() bool {
    return false
}
