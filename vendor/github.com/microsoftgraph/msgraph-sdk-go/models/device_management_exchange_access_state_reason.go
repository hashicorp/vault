package models
// Device Exchange Access State Reason.
type DeviceManagementExchangeAccessStateReason int

const (
    // No access state reason discovered from Exchange
    NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON DeviceManagementExchangeAccessStateReason = iota
    // Unknown access state reason
    UNKNOWN_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state determined by Exchange Global rule
    EXCHANGEGLOBALRULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state determined by Exchange Individual rule
    EXCHANGEINDIVIDUALRULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state determined by Exchange Device rule
    EXCHANGEDEVICERULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state due to Exchange upgrade
    EXCHANGEUPGRADE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state determined by Exchange Mailbox Policy
    EXCHANGEMAILBOXPOLICY_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state determined by Exchange
    OTHER_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state granted by compliance challenge
    COMPLIANT_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state revoked by compliance challenge
    NOTCOMPLIANT_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state revoked by management challenge
    NOTENROLLED_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state due to unknown location
    UNKNOWNLOCATION_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state due to MFA challenge
    MFAREQUIRED_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access State revoked by AAD Access Policy
    AZUREADBLOCKDUETOACCESSPOLICY_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access State revoked by compromised password
    COMPROMISEDPASSWORD_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    // Access state revoked by managed application challenge
    DEVICENOTKNOWNWITHMANAGEDAPP_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
)

func (i DeviceManagementExchangeAccessStateReason) String() string {
    return []string{"none", "unknown", "exchangeGlobalRule", "exchangeIndividualRule", "exchangeDeviceRule", "exchangeUpgrade", "exchangeMailboxPolicy", "other", "compliant", "notCompliant", "notEnrolled", "unknownLocation", "mfaRequired", "azureADBlockDueToAccessPolicy", "compromisedPassword", "deviceNotKnownWithManagedApp"}[i]
}
func ParseDeviceManagementExchangeAccessStateReason(v string) (any, error) {
    result := NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
    switch v {
        case "none":
            result = NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "unknown":
            result = UNKNOWN_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "exchangeGlobalRule":
            result = EXCHANGEGLOBALRULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "exchangeIndividualRule":
            result = EXCHANGEINDIVIDUALRULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "exchangeDeviceRule":
            result = EXCHANGEDEVICERULE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "exchangeUpgrade":
            result = EXCHANGEUPGRADE_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "exchangeMailboxPolicy":
            result = EXCHANGEMAILBOXPOLICY_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "other":
            result = OTHER_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "compliant":
            result = COMPLIANT_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "notCompliant":
            result = NOTCOMPLIANT_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "notEnrolled":
            result = NOTENROLLED_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "unknownLocation":
            result = UNKNOWNLOCATION_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "mfaRequired":
            result = MFAREQUIRED_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "azureADBlockDueToAccessPolicy":
            result = AZUREADBLOCKDUETOACCESSPOLICY_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "compromisedPassword":
            result = COMPROMISEDPASSWORD_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        case "deviceNotKnownWithManagedApp":
            result = DEVICENOTKNOWNWITHMANAGEDAPP_DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementExchangeAccessStateReason(values []DeviceManagementExchangeAccessStateReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementExchangeAccessStateReason) isMultiValue() bool {
    return false
}
