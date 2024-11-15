package models
// Delivery optimization mode for peer distribution
type WindowsDeliveryOptimizationMode int

const (
    // Allow the user to set.
    USERDEFINED_WINDOWSDELIVERYOPTIMIZATIONMODE WindowsDeliveryOptimizationMode = iota
    // HTTP only, no peering
    HTTPONLY_WINDOWSDELIVERYOPTIMIZATIONMODE
    // OS default – Http blended with peering behind the same network address translator
    HTTPWITHPEERINGNAT_WINDOWSDELIVERYOPTIMIZATIONMODE
    // HTTP blended with peering across a private group
    HTTPWITHPEERINGPRIVATEGROUP_WINDOWSDELIVERYOPTIMIZATIONMODE
    // HTTP blended with Internet peering
    HTTPWITHINTERNETPEERING_WINDOWSDELIVERYOPTIMIZATIONMODE
    // Simple download mode with no peering
    SIMPLEDOWNLOAD_WINDOWSDELIVERYOPTIMIZATIONMODE
    // Bypass mode. Do not use Delivery Optimization and use BITS instead
    BYPASSMODE_WINDOWSDELIVERYOPTIMIZATIONMODE
)

func (i WindowsDeliveryOptimizationMode) String() string {
    return []string{"userDefined", "httpOnly", "httpWithPeeringNat", "httpWithPeeringPrivateGroup", "httpWithInternetPeering", "simpleDownload", "bypassMode"}[i]
}
func ParseWindowsDeliveryOptimizationMode(v string) (any, error) {
    result := USERDEFINED_WINDOWSDELIVERYOPTIMIZATIONMODE
    switch v {
        case "userDefined":
            result = USERDEFINED_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "httpOnly":
            result = HTTPONLY_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "httpWithPeeringNat":
            result = HTTPWITHPEERINGNAT_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "httpWithPeeringPrivateGroup":
            result = HTTPWITHPEERINGPRIVATEGROUP_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "httpWithInternetPeering":
            result = HTTPWITHINTERNETPEERING_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "simpleDownload":
            result = SIMPLEDOWNLOAD_WINDOWSDELIVERYOPTIMIZATIONMODE
        case "bypassMode":
            result = BYPASSMODE_WINDOWSDELIVERYOPTIMIZATIONMODE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsDeliveryOptimizationMode(values []WindowsDeliveryOptimizationMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsDeliveryOptimizationMode) isMultiValue() bool {
    return false
}
