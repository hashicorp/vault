package models
type DelegateMeetingMessageDeliveryOptions int

const (
    SENDTODELEGATEANDINFORMATIONTOPRINCIPAL_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS DelegateMeetingMessageDeliveryOptions = iota
    SENDTODELEGATEANDPRINCIPAL_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
    SENDTODELEGATEONLY_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
)

func (i DelegateMeetingMessageDeliveryOptions) String() string {
    return []string{"sendToDelegateAndInformationToPrincipal", "sendToDelegateAndPrincipal", "sendToDelegateOnly"}[i]
}
func ParseDelegateMeetingMessageDeliveryOptions(v string) (any, error) {
    result := SENDTODELEGATEANDINFORMATIONTOPRINCIPAL_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
    switch v {
        case "sendToDelegateAndInformationToPrincipal":
            result = SENDTODELEGATEANDINFORMATIONTOPRINCIPAL_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
        case "sendToDelegateAndPrincipal":
            result = SENDTODELEGATEANDPRINCIPAL_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
        case "sendToDelegateOnly":
            result = SENDTODELEGATEONLY_DELEGATEMEETINGMESSAGEDELIVERYOPTIONS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDelegateMeetingMessageDeliveryOptions(values []DelegateMeetingMessageDeliveryOptions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DelegateMeetingMessageDeliveryOptions) isMultiValue() bool {
    return false
}
