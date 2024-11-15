package models
// Top level failure categories for enrollment.
type DeviceEnrollmentFailureReason int

const (
    // Default value, failure reason is unknown.
    UNKNOWN_DEVICEENROLLMENTFAILUREREASON DeviceEnrollmentFailureReason = iota
    // Authentication failed
    AUTHENTICATION_DEVICEENROLLMENTFAILUREREASON
    // Call was authenticated, but not authorized to enroll.
    AUTHORIZATION_DEVICEENROLLMENTFAILUREREASON
    // Failed to validate the account for enrollment. (Account blocked, enrollment not enabled)
    ACCOUNTVALIDATION_DEVICEENROLLMENTFAILUREREASON
    // User could not be validated. (User does not exist, missing license)
    USERVALIDATION_DEVICEENROLLMENTFAILUREREASON
    // Device is not supported for mobile device management.
    DEVICENOTSUPPORTED_DEVICEENROLLMENTFAILUREREASON
    // Account is in maintenance.
    INMAINTENANCE_DEVICEENROLLMENTFAILUREREASON
    // Client sent a request that is not understood/supported by the service.
    BADREQUEST_DEVICEENROLLMENTFAILUREREASON
    // Feature(s) used by this enrollment are not supported for this account.
    FEATURENOTSUPPORTED_DEVICEENROLLMENTFAILUREREASON
    // Enrollment restrictions configured by admin blocked this enrollment.
    ENROLLMENTRESTRICTIONSENFORCED_DEVICEENROLLMENTFAILUREREASON
    // Client timed out or enrollment was aborted by enduser.
    CLIENTDISCONNECTED_DEVICEENROLLMENTFAILUREREASON
    // Enrollment was abandoned by enduser. (Enduser started onboarding but failed to complete it in timely manner)
    USERABANDONMENT_DEVICEENROLLMENTFAILUREREASON
)

func (i DeviceEnrollmentFailureReason) String() string {
    return []string{"unknown", "authentication", "authorization", "accountValidation", "userValidation", "deviceNotSupported", "inMaintenance", "badRequest", "featureNotSupported", "enrollmentRestrictionsEnforced", "clientDisconnected", "userAbandonment"}[i]
}
func ParseDeviceEnrollmentFailureReason(v string) (any, error) {
    result := UNKNOWN_DEVICEENROLLMENTFAILUREREASON
    switch v {
        case "unknown":
            result = UNKNOWN_DEVICEENROLLMENTFAILUREREASON
        case "authentication":
            result = AUTHENTICATION_DEVICEENROLLMENTFAILUREREASON
        case "authorization":
            result = AUTHORIZATION_DEVICEENROLLMENTFAILUREREASON
        case "accountValidation":
            result = ACCOUNTVALIDATION_DEVICEENROLLMENTFAILUREREASON
        case "userValidation":
            result = USERVALIDATION_DEVICEENROLLMENTFAILUREREASON
        case "deviceNotSupported":
            result = DEVICENOTSUPPORTED_DEVICEENROLLMENTFAILUREREASON
        case "inMaintenance":
            result = INMAINTENANCE_DEVICEENROLLMENTFAILUREREASON
        case "badRequest":
            result = BADREQUEST_DEVICEENROLLMENTFAILUREREASON
        case "featureNotSupported":
            result = FEATURENOTSUPPORTED_DEVICEENROLLMENTFAILUREREASON
        case "enrollmentRestrictionsEnforced":
            result = ENROLLMENTRESTRICTIONSENFORCED_DEVICEENROLLMENTFAILUREREASON
        case "clientDisconnected":
            result = CLIENTDISCONNECTED_DEVICEENROLLMENTFAILUREREASON
        case "userAbandonment":
            result = USERABANDONMENT_DEVICEENROLLMENTFAILUREREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceEnrollmentFailureReason(values []DeviceEnrollmentFailureReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceEnrollmentFailureReason) isMultiValue() bool {
    return false
}
