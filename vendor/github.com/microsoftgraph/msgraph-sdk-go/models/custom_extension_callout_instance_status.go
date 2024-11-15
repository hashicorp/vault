package models
type CustomExtensionCalloutInstanceStatus int

const (
    CALLOUTSENT_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS CustomExtensionCalloutInstanceStatus = iota
    CALLBACKRECEIVED_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
    CALLOUTFAILED_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
    CALLBACKTIMEDOUT_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
    WAITINGFORCALLBACK_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
    UNKNOWNFUTUREVALUE_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
)

func (i CustomExtensionCalloutInstanceStatus) String() string {
    return []string{"calloutSent", "callbackReceived", "calloutFailed", "callbackTimedOut", "waitingForCallback", "unknownFutureValue"}[i]
}
func ParseCustomExtensionCalloutInstanceStatus(v string) (any, error) {
    result := CALLOUTSENT_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
    switch v {
        case "calloutSent":
            result = CALLOUTSENT_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        case "callbackReceived":
            result = CALLBACKRECEIVED_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        case "calloutFailed":
            result = CALLOUTFAILED_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        case "callbackTimedOut":
            result = CALLBACKTIMEDOUT_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        case "waitingForCallback":
            result = WAITINGFORCALLBACK_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CUSTOMEXTENSIONCALLOUTINSTANCESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCustomExtensionCalloutInstanceStatus(values []CustomExtensionCalloutInstanceStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CustomExtensionCalloutInstanceStatus) isMultiValue() bool {
    return false
}
